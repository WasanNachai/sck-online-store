import { IsNumber, Min } from 'class-validator';

export class CalculatePointDto {
@IsNumber()
@Min(0)
amountTHB: number;
}
