package utils

func GetRankEmote(rank string) string {
	switch rank {
	case "XH":
		return "<:rankingSSH:1107257232850174012>"
	case "X":
		return "<:rankingSS:1107257211866071050>"
	case "SH":
		return "<:rankingSH:1107257262860423238>"
	case "S":
		return "<:rankingS:1107257185974620321>"
	case "A":
		return "<:rankingA:1107256608544796702>"
	case "B":
		return "<:rankingB:1107257121407500288>"
	case "C":
		return "<:rankingC:1107257135194177647>"
	case "D":
		return "<:rankingD:1107257148867629120>"
	default:
		return "<:rankingSSH:1107257232850174012>"
	}
}
