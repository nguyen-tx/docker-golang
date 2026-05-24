// ===== User =====
export type UserRole = 'admin' | 'pilot' | 'operator' | 'atc' | 'inspector'
export type UserStatus = 'active' | 'inactive' | 'pending' | 'banned'

export interface User {
  id: string
  email: string
  full_name: string
  phone: string
  role: UserRole
  status: UserStatus
  license_no?: string
  organization?: string
}

// ===== Aircraft =====
export type AircraftType = 'fixed_wing' | 'multi_rotor' | 'helicopter' | 'hybrid'
export type AircraftStatus = 'active' | 'inactive' | 'maintenance' | 'grounded'

export interface Aircraft {
  id: string
  owner_id: string
  registration_no: string
  model: string
  manufacturer: string
  serial_no: string
  type: AircraftType
  status: AircraftStatus
  weight_kg: number
  max_altitude_m: number
  max_speed_kmh: number
  max_flight_time_min: number
  has_camera: boolean
  insurance_no?: string
  insurance_expiry?: string
  created_at: string
  owner?: User
}

// ===== Airspace =====
export type ZoneType = 'restricted' | 'controlled' | 'uncontrolled' | 'tfr' | 'nfz'

export interface AirspaceZone {
  id: string
  name: string
  description: string
  type: ZoneType
  min_altitude_m: number
  max_altitude_m: number
  boundary_geojson: string
  active_from?: string
  active_until?: string
  is_active: boolean
  created_at: string
}

// ===== Authorization =====
export type AuthorizationStatus = 'pending' | 'approved' | 'rejected' | 'cancelled' | 'expired' | 'revoked'
export type FlightPurpose = 'commercial' | 'recreational' | 'survey' | 'delivery' | 'emergency' | 'research'

export interface AuthorizationRequest {
  id: string
  request_no: string
  pilot_id: string
  aircraft_id: string
  zone_id?: string
  purpose: FlightPurpose
  status: AuthorizationStatus
  flight_area_json: string
  planned_alt_m: number
  planned_start_at: string
  planned_end_at: string
  pilot_notes?: string
  reviewer_id?: string
  reviewed_at?: string
  review_notes?: string
  valid_from?: string
  valid_until?: string
  created_at: string
  pilot?: User
  aircraft?: Aircraft
}

// ===== Telemetry =====
export interface Telemetry {
  id: string
  aircraft_id: string
  session_id: string
  auth_request_id?: string
  lat: number
  lng: number
  altitude_m: number
  speed_kmh: number
  heading_deg: number
  battery_pct: number
  signal_strength: number
  is_on_ground: boolean
  timestamp: string
}

export interface FlightSession {
  id: string
  session_id: string
  aircraft_id: string
  pilot_id: string
  auth_request_id?: string
  started_at: string
  ended_at?: string
  max_altitude_m: number
  total_dist_km: number
  is_active: boolean
}

// ===== Notification =====
export type NotificationType =
  | 'auth_approved' | 'auth_rejected' | 'auth_expiring'
  | 'zone_alert' | 'weather_alert' | 'system_info' | 'violation'

export type NotificationPriority = 'low' | 'medium' | 'high' | 'critical'

export interface Notification {
  id: string
  recipient_id: string
  type: NotificationType
  priority: NotificationPriority
  title: string
  body: string
  metadata?: Record<string, unknown>
  is_read: boolean
  read_at?: string
  created_at: string
}

// ===== API =====
export interface ApiError {
  error: string
}

export interface LoginResponse {
  access_token: string
  user: User
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}
