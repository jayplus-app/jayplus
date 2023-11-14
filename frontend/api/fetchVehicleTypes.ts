export interface VehicleType {
  id: number
  businessId: number
  name: string
  icon: string
  description: string
  position: number
}

export async function fetchVehicleTypes(): Promise<VehicleType[]> {
  try {
    const response = await fetch('/api/booking/vehicle-types')
    const data: VehicleType[] = await response.json()
    return data
  } catch (error) {
    console.error(error)
    throw new Error('Failed to fetch vehicle types data.')
  }
}
