import { Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Point } from './point.entity';
import { CreatePointDto } from './point.dto';
import { logs, SeverityNumber } from '@opentelemetry/api-logs';

const otelLogger = logs.getLogger('point-service');

export interface CalculatePointResponse {
  earnedPoints: number;
  amountTHB: number;
  rateTHBPerPoint: number;
}

@Injectable()
export class PointService {
  private readonly logger = new Logger(PointService.name);

  // New reward rule:
  // Every complete 50 THB = 1 reward point
  // Equivalent to 100 THB = 2 reward points
  private static readonly THB_PER_POINT = 50;

  constructor(
    @InjectRepository(Point)
    private readonly pointRepository: Repository<Point>,
  ) {}

  /**
   * Calculate reward points.
   */
  calculateEarnedPoints(amountTHB: number): number {
    if (!Number.isFinite(amountTHB) || amountTHB <= 0) {
      return 0;
    }

    return Math.floor(amountTHB / PointService.THB_PER_POINT);
  }

  /**
   * Response for reward point calculation API.
   */
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

  async getPoint(): Promise<Point[]> {
    try {
      const points = await this.pointRepository.find();

      this.logger.log(`Points retrieved, count=${points.length}`);

      otelLogger.emit({
        severityNumber: SeverityNumber.INFO,
        severityText: 'INFO',
        body: 'Points retrieved',
        attributes: {
          log_type: 'business',
          event: 'points_retrieved',
          entity_type: 'point',
          items_count: points.length,
        },
      });

      return points;
    } catch (error: any) {
      this.logger.error(
        'PointRepository.find internal error',
        error.stack,
      );

      otelLogger.emit({
        severityNumber: SeverityNumber.ERROR,
        severityText: 'ERROR',
        body: 'PointRepository.find internal error',
        attributes: {
          'error.message': error.message,
        },
      });

      throw error;
    }
  }

  async deductPoint(point: CreatePointDto): Promise<Point> {
    try {
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
    } catch (error: any) {
      this.logger.error(
        'PointRepository.save internal error',
        error.stack,
      );

      otelLogger.emit({
        severityNumber: SeverityNumber.ERROR,
        severityText: 'ERROR',
        body: 'PointRepository.save internal error',
        attributes: {
          'error.message': error.message,
        },
      });

      throw error;
    }
  }
}