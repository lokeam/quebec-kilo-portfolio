import type { BaseLocation, PhysicalLocation, Sublocation } from '../types/location';
import { axiosInstance } from '../client/axios-instance';
import { isPhysicalLocation, isSublocation } from '../types/location';
import type { AxiosResponse } from 'axios';
import { toCamelCase } from '../utils/serialization';

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

  public async createPhysicalLocation(location: PhysicalLocation): Promise<PhysicalLocation> {
    const response: AxiosResponse<PhysicalLocation> = await axiosInstance.post('/v1/locations/physical', location);
    return toCamelCase(response.data);
  }

  public async createSublocation(location: Sublocation): Promise<Sublocation> {
    const response: AxiosResponse<Sublocation> = await axiosInstance.post('/v1/locations/sublocations', location);
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

  public async updatePhysicalLocation(location: PhysicalLocation): Promise<PhysicalLocation> {
    const response: AxiosResponse<PhysicalLocation> = await axiosInstance.put(`/v1/locations/physical/${location.id}`, location);
    return toCamelCase(response.data);
  }

  public async updateSublocation(location: Sublocation): Promise<Sublocation> {
    // Convert camelCase to snake_case for API request
    const apiPayload = {
      name: location.name,
      location_type: location.locationType,
      map_coordinates: location.mapCoordinates,
      bg_color: location.bgColor,
      physical_location_id: location.parentLocationId,
    };

    const response: AxiosResponse<Sublocation> = await axiosInstance.put(
      `/v1/locations/sublocations/${location.id}`,
      apiPayload
    );
    return toCamelCase(response.data);
  }

  async deleteLocation(id: string, type: 'physical' | 'sublocation'): Promise<void> {
    const endpoint = type === 'physical' ? '/v1/locations/physical' : '/v1/locations/sublocations';
    await axiosInstance.delete(`${endpoint}/${id}`);
  }

  public async getPhysicalLocations(): Promise<PhysicalLocation[]> {
    const response: AxiosResponse<{ locations: PhysicalLocation[] }> = await axiosInstance.get('/v1/locations/physical');
    return toCamelCase(response.data.locations);
  }

  public async getSublocations(): Promise<Sublocation[]> {
    const response: AxiosResponse<{ locations: Sublocation[] }> = await axiosInstance.get('/v1/locations/sublocations');
    return toCamelCase(response.data.locations);
  }
}