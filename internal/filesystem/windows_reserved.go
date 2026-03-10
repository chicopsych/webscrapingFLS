package filesystem

// windowsReservedNames representa nomes de dispositivo reservados pelo namespace
// Win32/NTFS. Esses identificadores não podem ser usados como nome de arquivo,
// mesmo com extensão (ex.: CON.txt continua inválido).
//
// Referência histórica: DOS devices mapeados pela API do Windows.
var windowsReservedNames = map[string]struct{}{
	"CON":  {},
	"PRN":  {},
	"AUX":  {},
	"NUL":  {},
	"COM1": {},
	"COM2": {},
	"COM3": {},
	"COM4": {},
	"COM5": {},
	"COM6": {},
	"COM7": {},
	"COM8": {},
	"COM9": {},
	"LPT1": {},
	"LPT2": {},
	"LPT3": {},
	"LPT4": {},
	"LPT5": {},
	"LPT6": {},
	"LPT7": {},
	"LPT8": {},
	"LPT9": {},
}
