import {
  Body,
  Controller,
  Get,
  HttpException,
  HttpStatus,
  Logger,
  Post,
} from '@nestjs/common';
import { PointService } from './point.service';
import { CreatePointDto } from './point.dto';
import { CalculatePointDto } from './calculate-point.dto';
import { logs, SeverityNumber } from '@opentelemetry/api-logs';

const otelLogger = logs.getLogger('point-service');

@Controller('point')
export class PointController {
  private readonly logger = new Logger(PointController.name);

  constructor(private readonly pointService: PointService) {}

  @Get()
  async getPoint() {
    this.logger.log('GET /point request received');

    otelLogger.emit({
      severityNumber: SeverityNumber.INFO,
      severityText: 'INFO',
      body: 'Get points request received',
      attributes: {
        log_type: 'business',
        event: 'get_points_request',
        entity_type: 'point',
      },
    });

    try {
      return await this.pointService.getPoint();
    } catch (error: any) {
      this.logger.error(
        'PointService.getPoint internal error',
        error.stack,
      );

      otelLogger.emit({
        severityNumber: SeverityNumber.ERROR,
        severityText: 'ERROR',
        body: 'PointService.getPoint internal error',
        attributes: {
          'error.message': error.message,
        },
      });

      throw new HttpException(
        error.message,
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  @Post()
  async createPoint(@Body() body: CreatePointDto) {
    this.logger.log(
      `POST /point request received: userId=${body.userId}, orgId=${body.orgId}, amount=${body.amount}`,
    );

    otelLogger.emit({
      severityNumber: SeverityNumber.INFO,
      severityText: 'INFO',
      body: 'Deduct points request received',
      attributes: {
        log_type: 'business',
        event: 'deduct_points_request',
        entity_type: 'point',
        actor_id: body.userId,
        org_id: body.orgId,
        amount: body.amount,
      },
    });

    try {
      return await this.pointService.deductPoint(body);
    } catch (error: any) {
      this.logger.error(
        'PointService.deductPoint internal error',
        error.stack,
      );

      otelLogger.emit({
        severityNumber: SeverityNumber.ERROR,
        severityText: 'ERROR',
        body: 'PointService.deductPoint internal error',
        attributes: {
          'error.message': error.message,
        },
      });

      throw new HttpException(
        error.message,
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  @Post('calculate')
  calculatePoint(@Body() body: CalculatePointDto) {
    this.logger.log(
      `POST /point/calculate request received: amountTHB=${body.amountTHB}`,
    );

    otelLogger.emit({
      severityNumber: SeverityNumber.INFO,
      severityText: 'INFO',
      body: 'Calculate reward points request received',
      attributes: {
        log_type: 'business',
        event: 'calculate_reward_points_request',
        entity_type: 'point',
        amount_thb: body.amountTHB,
      },
    });

    try {
      return this.pointService.calculatePointResponse(body.amountTHB);
    } catch (error: any) {
      this.logger.error(
        'PointService.calculatePointResponse internal error',
        error.stack,
      );

      otelLogger.emit({
        severityNumber: SeverityNumber.ERROR,
        severityText: 'ERROR',
        body: 'PointService.calculatePointResponse internal error',
        attributes: {
          'error.message': error.message,
          amount_thb: body.amountTHB,
        },
      });

      throw new HttpException(
        error.message,
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
}