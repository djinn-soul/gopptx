package text

import "strings"

const (
	// MaxBulletLevel defines the maximum nested bullet depth supported by PowerPoint.
	MaxBulletLevel = 9

	// BulletStyleBullet renders a standard bullet character.
	BulletStyleBullet = "bullet"
	// BulletStyleNumber renders arabic numbering (1., 2., 3.).
	BulletStyleNumber = "number"
	// BulletStyleLetterLower renders lowercase lettering (a., b., c.).
	BulletStyleLetterLower = "letter_lower"
	// BulletStyleLetterUpper renders uppercase lettering (A., B., C.).
	BulletStyleLetterUpper = "letter_upper"
	// BulletStyleRomanLower renders lowercase roman numerals (i., ii., iii.).
	BulletStyleRomanLower = "roman_lower"
	// BulletStyleRomanUpper renders uppercase roman numerals (I., II., III.).
	BulletStyleRomanUpper = "roman_upper"
	// BulletStyleCustom renders one caller-provided bullet character.
	BulletStyleCustom = "custom"
	// BulletStyleNone renders no bullet marker.
	BulletStyleNone = "none"
)

// NormalizeBulletStyle sanitizes bullet style strings.
func NormalizeBulletStyle(style string) string {
	normalized := strings.ToLower(strings.TrimSpace(style))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")

	switch normalized {
	case "":
		return ""
	case BulletStyleBullet:
		return BulletStyleBullet
	case BulletStyleNumber, "numbered":
		return BulletStyleNumber
	case BulletStyleLetterLower, "lettered", "letter", "letterlower", "alphalower":
		return BulletStyleLetterLower
	case BulletStyleLetterUpper, "letterupper", "alphaupper":
		return BulletStyleLetterUpper
	case BulletStyleRomanLower, "romanlower":
		return BulletStyleRomanLower
	case BulletStyleRomanUpper, "roman", "romanupper":
		return BulletStyleRomanUpper
	case BulletStyleCustom:
		return BulletStyleCustom
	case BulletStyleNone:
		return BulletStyleNone
	default:
		return normalized
	}
}
