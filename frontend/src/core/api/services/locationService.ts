import type { BaseLocation, PhysicalLocation, Sublocation } from '../types/location';
import { axiosInstance } from '../client/axios-instance';
import { isPhysicalLocation, isSublocation } from '../types/location';
import type { AxiosResponse } from 'axios';
import { toCamelCase, toSnakeCase } from '../utils/serialization';

export class LocationService {
  private static instance: LocationService;

  private constructor() {}

  public static getInstance(): LocationService {
    if (!LocationService.instance) {
      LocationService.instance = new LocationService();
    }
    return LocationService.instance;
  }

  async createLocation(location: BaseLocation) {
    if (isPhysicalLocation(location)) {
      return this.createPhysicalLocation(location);
    } else if (isSublocation(location)) {
      return this.createSublocation(location);
    }
    throw new Error('Invalid location type');
  }

  private async createPhysicalLocation(location: PhysicalLocation): Promise<PhysicalLocation> {
    const snakeCaseLocation = toSnakeCase(location);
    const response: AxiosResponse<PhysicalLocation> = await axiosInstance.post('/v1/locations/physical', snakeCaseLocation);
    return toCamelCase(response.data);
  }

  private async createSublocation(location: Sublocation): Promise<Sublocation> {
    const snakeCaseLocation = toSnakeCase(location);
    const response: AxiosResponse<Sublocation> = await axiosInstance.post('/v1/locations/sublocations', snakeCaseLocation);
    return toCamelCase(response.data);
  }

  async updateLocation(location: BaseLocation) {
    if (isPhysicalLocation(location)) {
      return this.updatePhysicalLocation(location);
    } else if (isSublocation(location)) {
      return this.updateSublocation(location);
    }
    throw new Error('Invalid location type');
  }

  private async updatePhysicalLocation(location: PhysicalLocation): Promise<PhysicalLocation> {
    const snakeCaseLocation = toSnakeCase(location);
    const response: AxiosResponse<PhysicalLocation> = await axiosInstance.put(`/v1/locations/physical/${location.id}`, snakeCaseLocation);
    return toCamelCase(response.data);
  }

  private async updateSublocation(location: Sublocation): Promise<Sublocation> {
    const snakeCaseLocation = toSnakeCase(location);
    const response: AxiosResponse<Sublocation> = await axiosInstance.put(`/v1/locations/sublocations/${location.id}`, snakeCaseLocation);
    return toCamelCase(response.data);
  }

  async deleteLocation(id: string, type: 'physical' | 'sublocation'): Promise<void> {
    const endpoint = type === 'physical' ? '/v1/locations/physical' : '/v1/locations/sublocations';
    await axiosInstance.delete(`${endpoint}/${id}`);
  }
}