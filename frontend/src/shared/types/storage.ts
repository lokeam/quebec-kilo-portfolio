/**
 * Units of digital storage measurement
 * Ordered from smallest to largest size
 *
 * - B: Bytes
 * - KB: Kilobytes (1024 bytes)
 * - MB: Megabytes (1024 KB)
 * - GB: Gigabytes (1024 MB)
 * - TB: Terabytes (1024 GB)
 */
export type StorageSizeUnit =
  | 'B'
  | 'KB'
  | 'MB'
  | 'GB'
  | 'TB';

/**
 * Represents a digital storage size with a value + unit
 *
 * @example
 * ```ts
 * // 500 GB hard drive
 * const hddSize: StorageSize = {
 *   value: "500",
 *   unit: "GB"
 * };
 *
 * // 2.5 MB file
 * const fileSize: StorageSize = {
 *   value: "2.5",
 *   unit: "MB"
 * };
 * ```
 */
export interface StorageSize {
   /** Numeric value as a string to support decimal points */
  value: string;
  /** Unit of storage measurement */
  unit: StorageSizeUnit;
}