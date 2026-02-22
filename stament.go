package obreron

// NOTE: The core internal engine (stament/segment/scanner/build) is implemented
// in subsequent issues. This file exists to anchor the internal API and keep
// the package structure stable.

import "github.com/profe-ajedrez/obreron/v3/dialect"

type stament struct {
	buf       []byte          // texto de cláusulas con marcadores 0x1F para placeholders
	segs      []segment       // descriptores ordenables por sType
	params    []any           // parámetros posicionales en orden de inserción
	dialect   dialect.Dialect // inyectado desde DB o constructor directo
	err       error           // primer error acumulado — subsiguientes se descartan
	lastSType uint8           // último sType, para el método Clause()
	flags     uint8           // bitfield de estado (ver constantes abajo)
}

const (
	flagGrouped    uint8 = 1 << 0
	flagFirstCol   uint8 = 1 << 1
	flagWhereAdded uint8 = 1 << 2
	flagWithSelect uint8 = 1 << 3
	flagHasFrom    uint8 = 1 << 4
	flagHasTable   uint8 = 1 << 5
	flagHasSet     uint8 = 1 << 6
	flagHasCols    uint8 = 1 << 7
)

const (
	selectKW uint8 = 0  // keyword inicial: SELECT / DELETE / UPDATE / INSERT
	colsS    uint8 = 1  // columnas SELECT o lista de columnas INSERT
	setS     uint8 = 2  // SET (UPDATE)
	fromS    uint8 = 3  // FROM
	joinS    uint8 = 4  // JOIN / LEFT JOIN / RIGHT JOIN / CROSS JOIN
	whereS   uint8 = 5  // WHERE / AND / OR
	groupS   uint8 = 6  // GROUP BY
	havingS  uint8 = 7  // HAVING
	orderS   uint8 = 8  // ORDER BY
	limitS   uint8 = 9  // LIMIT
	offsetS  uint8 = 10 // OFFSET
	retS     uint8 = 11 // RETURNING
	insValS  uint8 = 99 // VALUES / SELECT subquery (INSERT)
)

func resetStament(st *stament) {
	clear(st.params)
	st.params = st.params[:0]
	clear(st.segs)
	st.segs = st.segs[:0]
	st.buf = st.buf[:0]
	st.dialect = nil // ← obligatorio: evita retener refs al pool
	st.err = nil
	st.lastSType = 0
	st.flags = flagFirstCol // único flag activo por defecto
}

// Layout exacto — 16 bytes en amd64:
//
//	start   uint32  →  4 bytes
//	length  uint32  →  4 bytes
//	pIndex  int32   →  4 bytes
//	pCount  uint16  →  2 bytes
//	sType   uint8   →  1 byte
//	_pad    uint8   →  1 byte  (padding explícito)
type segment struct {
	start  uint32 // offset en buf. Máx: 4 GB
	length uint32 // largo en buf. Máx: 4 GB
	pIndex int32  // índice base en params[]. -1 si no hay params
	pCount uint16 // cantidad de params de esta cláusula. Máx: 65535
	sType  uint8  // tipo / orden en el SQL final
	_pad   uint8  // padding explícito — reservado
}
