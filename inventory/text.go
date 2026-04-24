package inventory

import (
	"fmt"
	"io"
	"strings"
)

const (
	emDash             = "—"
	textMinPath        = 20
	textMinType        = 10
	textMinKey         = 20
	textMaxValueLength = 40
)

// WriteText writes a human-readable, fixed-width-column inventory to writer.
// Columns: PATH | TYPE | KEY | DEFAULT. Empty cells render as "—".
// Long default values are truncated with ellipsis. Output ends with a
// trailing newline.
func (r *Report) WriteText(writer io.Writer) error {
	const errCtx = "writing text inventory"

	rows := buildTextRows(r)

	pathWidth := columnWidth(rows, 0, textMinPath)
	typeWidth := columnWidth(rows, 1, textMinType)
	keyWidth := columnWidth(rows, 2, textMinKey)

	var buf strings.Builder
	fmt.Fprintf(&buf, "TYPE: %s\n\n", r.Type)
	fmt.Fprintf(
		&buf,
		"%-*s  %-*s  %-*s  %s\n",
		pathWidth, "PATH",
		typeWidth, "TYPE",
		keyWidth, "KEY",
		"DEFAULT",
	)

	for _, row := range rows {
		fmt.Fprintf(
			&buf,
			"%-*s  %-*s  %-*s  %s\n",
			pathWidth, row[0],
			typeWidth, row[1],
			keyWidth, row[2],
			row[3],
		)
	}

	if _, err := io.WriteString(writer, buf.String()); err != nil {
		return fmt.Errorf("%s: %w", errCtx, err)
	}

	return nil
}

// buildTextRows converts r.Fields into the [path, type, key, default]
// quadruples used by WriteText.
func buildTextRows(rep *Report) [][4]string {
	rows := make([][4]string, 0, len(rep.Fields))

	for _, fld := range rep.Fields {
		rows = append(rows, [4]string{
			fld.Path,
			fld.GoType,
			renderTextKey(fld.Key),
			renderTextDefault(fld.Satisfied),
		})
	}

	return rows
}

// renderTextKey renders a KeySpec as "<layer>: <key>" or em-dash.
func renderTextKey(ks *KeySpec) string {
	if ks == nil {
		return emDash
	}

	return ks.Layer + ": " + ks.Key
}

// renderTextDefault renders a Satisfaction as "<layerID>=<value>" with
// truncation, or em-dash when nil.
func renderTextDefault(sat *Satisfaction) string {
	if sat == nil {
		return emDash
	}

	val := fmt.Sprintf("%v", normalizeValue(sat.Value))
	if len(val) > textMaxValueLength {
		val = val[:textMaxValueLength-1] + "…"
	}

	return sat.LayerID + "=" + val
}

// columnWidth returns max(len(rows[col]), minimum) — the width needed to
// fit every value in the column.
func columnWidth(rows [][4]string, col, minimum int) int {
	width := minimum

	for _, row := range rows {
		if len(row[col]) > width {
			width = len(row[col])
		}
	}

	return width
}
