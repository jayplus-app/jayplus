import ServiceTypeSelection from './ServiceTypeSelection'
import TimeSelection from './TimeSelection'
import VehicleTypeSelection from './VehicleTypeSelection'

export default function Page() {
  return (
    <div>
      <h1>Booking Selection</h1>
      <VehicleTypeSelection />
      <ServiceTypeSelection />
      <TimeSelection />
    </div>
  )
}
