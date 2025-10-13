# Sistema de Gestión de Gimnasio - Diseño de Datos

## 📋 Tabla de Contenidos
- [Entidades de Base de Datos](#entidades-de-base-de-datos)
- [DTOs Públicos](#dtos-públicos-api--cliente)
- [DTOs Internos](#dtos-internos-entre-serviciosapi-interna)
- [DTOs de Request](#dtos-de-request-cliente--api)
- [Constantes y Enums](#constantes-y-enums)

---

## 🗄️ Entidades de Base de Datos

### GORM Entities

#### User
```go
type User struct {
    ID            uint      // Primary Key
    Email         string    // unique, not null
    Nombre        string    // not null
    Apellido      string    // not null
    FechaNac      time.Time // not null
    IsAdmin       bool      // default: false
    Telefono      string
    FechaRegistro time.Time // autoCreateTime
    Activo        bool      // default: true
}
```

**Índices:**
- `email` (unique)

**Validaciones:**
- Email debe ser único y válido
- Nombre y apellido obligatorios
- Fecha de nacimiento obligatoria (validar mayoría de edad)

---

#### Profesor
```go
type Profesor struct {
    ID              uint   // Primary Key
    UserID          uint   // Foreign Key -> User (unique)
    Especialidad    string
    Certificaciones string
}
```

**Relaciones:**
- `User` (1:1) - Un profesor está asociado a un usuario

**Notas:**
- No todos los usuarios son profesores
- Un usuario solo puede ser profesor una vez

---

### MongoDB Documents

#### Clase
```go
type Clase struct {
    ID                string    // _id (ObjectID)
    Nombre            string    // ej: "Yoga Principiantes"
    Descripcion       string
    ProfesorID        uint      // FK a Profesor (PostgreSQL)
    DiaSemana         string    // "Lunes", "Martes", etc
    HoraInicio        string    // "09:00"
    HoraFin           string    // "10:30"
    UsuariosInscritos []uint    // Array de User IDs
    CapacidadMax      int
    Activa            bool
    FechaCreacion     time.Time
}
```

**Validaciones:**
- `CapacidadMax` > 0
- `DiaSemana` debe estar en lista válida
- `HoraInicio` < `HoraFin`
- `len(UsuariosInscritos)` <= `CapacidadMax`

**Índices recomendados:**
- `profesor_id`
- `dia_semana` + `hora_inicio`
- `activa`

---

#### Asistencia
```go
type Asistencia struct {
    ID            string    // _id (ObjectID)
    ClaseID       string    // FK a Clase (MongoDB)
    UserID        uint      // FK a User (PostgreSQL)
    Fecha         time.Time // Fecha específica de la clase
    Asistio       bool
    Observaciones string    // opcional
}
```

**Índices recomendados:**
- `clase_id` + `user_id` + `fecha` (unique compound)
- `user_id` + `fecha`

---

## 🌐 DTOs Públicos (API -> Cliente)

Estos DTOs se usan para respuestas al cliente. **No incluyen información sensible**.

#### UserPublicDTO
```go
type UserPublicDTO struct {
    ID       uint   `json:"id"`
    Nombre   string `json:"nombre"`
    Apellido string `json:"apellido"`
    Email    string `json:"email"`
    Telefono string `json:"telefono,omitempty"`
}
```

---

#### ProfesorPublicDTO
```go
type ProfesorPublicDTO struct {
    ID           uint   `json:"id"`
    Nombre       string `json:"nombre"`
    Apellido     string `json:"apellido"`
    Especialidad string `json:"especialidad"`
}
```

---

#### ClasePublicDTO
```go
type ClasePublicDTO struct {
    ID                 string            `json:"id"`
    Nombre             string            `json:"nombre"`
    Descripcion        string            `json:"descripcion"`
    Profesor           ProfesorPublicDTO `json:"profesor"`
    DiaSemana          string            `json:"dia_semana"`
    HoraInicio         string            `json:"hora_inicio"`
    HoraFin            string            `json:"hora_fin"`
    LugaresDisponibles int               `json:"lugares_disponibles"`
    CapacidadMax       int               `json:"capacidad_max"`
}
```

**Nota:** `LugaresDisponibles` = `CapacidadMax - len(UsuariosInscritos)`

---

#### ClaseDetalleDTO
```go
type ClaseDetalleDTO struct {
    ClasePublicDTO           // Hereda todos los campos
    Inscritos []UserPublicDTO `json:"inscritos"`
}
```

**Uso:** Para vistas detalladas de una clase (ej: panel de admin)

---

## 🔒 DTOs Internos (Entre servicios/API interna)

Estos DTOs se usan para comunicación entre microservicios o capas internas.

#### UserInternoDTO
```go
type UserInternoDTO struct {
    UserID  uint   `json:"user_id"`
    Email   string `json:"email"`
    IsAdmin bool   `json:"is_admin"`
    Activo  bool   `json:"activo"`
}
```

**Uso:** Autenticación, autorización, eventos internos

---

#### ClaseInternaDTO
```go
type ClaseInternaDTO struct {
    ID                string `json:"id"`
    ProfesorID        uint   `json:"profesor_id"`
    UsuariosInscritos []uint `json:"usuarios_inscritos"`
    CapacidadMax      int    `json:"capacidad_max"`
    Activa            bool   `json:"activa"`
}
```

**Uso:** Validaciones de negocio, procesamiento interno

---

## 📥 DTOs de Request (Cliente -> API)

Estos DTOs validan datos de entrada del cliente.

#### CrearClaseRequest
```go
type CrearClaseRequest struct {
    Nombre       string `json:"nombre" binding:"required"`
    Descripcion  string `json:"descripcion"`
    ProfesorID   uint   `json:"profesor_id" binding:"required"`
    DiaSemana    string `json:"dia_semana" binding:"required"`
    HoraInicio   string `json:"hora_inicio" binding:"required"`
    HoraFin      string `json:"hora_fin" binding:"required"`
    CapacidadMax int    `json:"capacidad_max" binding:"required,min=1"`
}
```

**Validaciones adicionales:**
- `DiaSemana` debe estar en lista válida
- `HoraInicio` formato HH:MM
- `HoraFin` > `HoraInicio`
- `ProfesorID` debe existir en BD

---

#### InscribirseClaseRequest
```go
type InscribirseClaseRequest struct {
    UserID  uint   `json:"user_id" binding:"required"`
    ClaseID string `json:"clase_id" binding:"required"`
}
```

**Validaciones de negocio:**
- Usuario existe y está activo
- Clase existe y está activa
- Hay lugares disponibles
- Usuario no está ya inscrito

---

#### RegistrarAsistenciaRequest
```go
type RegistrarAsistenciaRequest struct {
    UserID  uint   `json:"user_id" binding:"required"`
    ClaseID string `json:"clase_id" binding:"required"`
    Asistio bool   `json:"asistio"`
}
```

**Validaciones de negocio:**
- Usuario está inscrito en la clase
- Fecha de registro válida (no futura)

---

## 📌 Constantes y Enums

### Días de la Semana
```go
var DiasSemana = []string{
    "Lunes", 
    "Martes", 
    "Miércoles", 
    "Jueves", 
    "Viernes", 
    "Sábado", 
    "Domingo",
}
```

### Roles (opcional - para futura expansión)
```go
const (
    RolUsuario   = "usuario"
    RolProfesor  = "profesor"
    RolAdmin     = "admin"
)
```