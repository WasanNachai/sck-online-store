import { PointStatus } from './point.entity';

export class CreatePointDto {
  orgId: number;
  userId: number;
  amount: number;
  status?: PointStatus;
  expireDate?: string;
}

export class CreatePendingEarnPointDto {
  orgId: number;
  userId: number;
  amount: number;
  expireDate: string;
}

export class TransitionPointDto {
  userId: number;
}

export class ExpirePointDto {
  beforeDate: string;
  userId?: number;
}
