import { Injectable, Logger, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Point, PointStatus } from './point.entity';
import {
  CreatePendingEarnPointDto,
  CreatePointDto,
} from './point.dto';
import { logs, SeverityNumber } from '@opentelemetry/api-logs';

const otelLogger = logs.getLogger('point-service');

export interface CalculatePointResponse {
  earnedPoints: number;
  amountTHB: number;
  rateTHBPerPoint: number;
}

export interface PointSummary {
  availablePoints: number;
  pendingPoints: number;
  redeemedPoints: number;
  expiredPoints: number;
}

@Injectable()
export class PointService {
  private readonly logger = new Logger(PointService.name);
  private static readonly THB_PER_POINT = 50;

  constructor(
    @InjectRepository(Point)
    private readonly pointRepository: Repository<Point>,
  ) {}

  calculateEarnedPoints(amountTHB: number): number {
    if (!Number.isFinite(amountTHB) || amountTHB <= 0) {
      return 0;
    }

    return Math.floor(amountTHB / PointService.THB_PER_POINT);
  }

  calculatePointResponse(
    amountTHB: number,
  ): CalculatePointResponse {
    const earnedPoints = this.calculateEarnedPoints(amountTHB);

    this.logger.log(
      `Reward points calculated: amountTHB=${amountTHB}, earnedPoints=${earnedPoints}`,
    );

    otelLogger.emit({
      severityNumber: SeverityNumber.INFO,
      severityText: 'INFO',
      body: 'Reward points calculated',
      attributes: {
        log_type: 'business',
        event: 'reward_points_calculated',
        entity_type: 'point',
        amount_thb: amountTHB,
        earned_points: earnedPoints,
        rate_thb_per_point: PointService.THB_PER_POINT,
      },
    });

    return {
      earnedPoints,
      amountTHB,
      rateTHBPerPoint: PointService.THB_PER_POINT,
    };
  }

  async getPoint(userId?: number): Promise<Point[]> {
    const points = await this.pointRepository.find({
      where: userId ? { userId } : {},
    });

    this.logger.log(`Points retrieved, count=${points.length}`);
    return points;
  }

  async deductPoint(point: CreatePointDto): Promise<Point> {
    const saved = await this.pointRepository.save(point);

    this.logger.log(
      `Points deducted: userId=${point.userId}, orgId=${point.orgId}, amount=${point.amount}`,
    );

    otelLogger.emit({
      severityNumber: SeverityNumber.INFO,
      severityText: 'INFO',
      body: 'Points deducted',
      attributes: {
        log_type: 'state_change',
        event: 'points_deducted',
        entity_type: 'point',
        entity_id: saved.id,
        changed_by: point.userId,
        org_id: point.orgId,
        amount: point.amount,
      },
    });

    return saved;
  }

  async createPendingEarnPoint(
    point: CreatePendingEarnPointDto,
  ): Promise<Point> {
    return this.pointRepository.save({
      orgId: point.orgId,
      userId: point.userId,
      amount: point.amount,
      status: PointStatus.PENDING_APPROVAL,
      expireDate: new Date(point.expireDate),
    });
  }

  async approvePoint(
    pointId: number,
    userId: number,
  ): Promise<Point> {
    const point = await this.pointRepository.findOne({
      where: {
        id: pointId,
        userId,
        status: PointStatus.PENDING_APPROVAL,
      },
    });

    if (!point) {
      throw new NotFoundException(
        'pending point not found for approval',
      );
    }

    point.status = PointStatus.APPROVED;
    return this.pointRepository.save(point);
  }

  async redeemPoint(
    pointId: number,
    userId: number,
  ): Promise<Point> {
    const point = await this.pointRepository.findOne({
      where: {
        id: pointId,
        userId,
        status: PointStatus.APPROVED,
      },
    });

    if (!point) {
      throw new NotFoundException(
        'approved point not found for redeem',
      );
    }

    point.status = PointStatus.REDEEMED;
    return this.pointRepository.save(point);
  }

  async expirePoints(
    beforeDate: Date,
    userId?: number,
  ): Promise<number> {
    const query = this.pointRepository
      .createQueryBuilder()
      .update(Point)
      .set({
        status: PointStatus.EXPIRED,
      })
      .where('status = :status', {
        status: PointStatus.APPROVED,
      })
      .andWhere('expire_date IS NOT NULL')
      .andWhere('expire_date < :beforeDate', {
        beforeDate,
      });

    if (userId) {
      query.andWhere('user_id = :userId', { userId });
    }

    const result = await query.execute();
    return result.affected ?? 0;
  }

  async getPointSummary(userId: number): Promise<PointSummary> {
    const points = await this.pointRepository.find({
      where: { userId },
    });
    const now = new Date();

    return points.reduce(
      (summary, point) => {
        switch (point.status) {
          case PointStatus.PENDING_APPROVAL:
            summary.pendingPoints += point.amount;
            break;
          case PointStatus.APPROVED:
            if (
              point.expireDate &&
              point.expireDate.getTime() < now.getTime()
            ) {
              summary.expiredPoints += point.amount;
            } else {
              summary.availablePoints += point.amount;
            }
            break;
          case PointStatus.REDEEMED:
            summary.redeemedPoints += point.amount;
            break;
          case PointStatus.EXPIRED:
            summary.expiredPoints += point.amount;
            break;
          default:
            break;
        }

        return summary;
      },
      {
        availablePoints: 0,
        pendingPoints: 0,
        redeemedPoints: 0,
        expiredPoints: 0,
      } as PointSummary,
    );
  }
}
