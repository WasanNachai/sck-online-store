import {
  Entity,
  Column,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  UpdateDateColumn,
} from 'typeorm';

export enum PointStatus {
  PENDING_APPROVAL = 'pending_approval',
  APPROVED = 'approved',
  REDEEMED = 'redeemed',
  EXPIRED = 'expired',
}

@Entity('points')
export class Point {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ name: 'org_id' })
  orgId: number;

  @Column({ name: 'user_id' })
  userId: number;

  @Column()
  amount: number;

  @Column({
    type: 'enum',
    enum: PointStatus,
    default: PointStatus.APPROVED,
  })
  status: PointStatus;

  @Column({
    name: 'expire_date',
    type: 'datetime',
    nullable: true,
  })
  expireDate: Date | null;

  @CreateDateColumn()
  created: Date;

  @UpdateDateColumn()
  updated: Date;
}
