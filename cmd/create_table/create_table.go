package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

// To align with iconv
var replaceMap = map[rune]rune{
	0xA1BD: 0x2015,
	0xA1B1: 0xFFE3,
	0xA1EF: 0xFFE5,
}

func main() {
	url := "http://x0213.org/codetable/euc-jis-2004-std.txt"
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer res.Body.Close()

	/*
	   0x00 - 0x7F ASCII
	   0xA1A1 - 0xFEFE 第一字面
	   0x8EA1 - 0x8EDF 半角かたかな
	   0x8FA1A1 - 0x8FFEFE 第二字面
	*/

	var eucJis1st = [23902]rune{}              //  A1A1(41377) ~ FEFE(65278)
	var eucJisHankaku = [63]rune{}             // 8EA1 (36513) ~  8EDF (36575)
	var eucJis2nd = [23894]rune{}              // 8FA1A1 (9413025) ~ 8FFEF6 (9436918)
	var eusJis1stCombined = map[rune][2]rune{} // 0xA1A1 ~ 0xFEFE
	var eusJis1stCombinedKeys = []rune{}

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "##") || strings.Contains(s, "# <reserved>") {
			continue
		}

		var euc_code int32
		var uni_code1 int32
		var uni_code2 int32

		parseLine(s, &euc_code, &uni_code1, &uni_code2)

		switch {
		case euc_code <= 0x7E:
			break

		case euc_code <= 0x9F:
			break

		case 0xA1A1 <= euc_code && euc_code <= 0xFEFE:
			if uni_code2 != 0 {
				eusJis1stCombinedKeys = append(eusJis1stCombinedKeys, rune(euc_code))
				eusJis1stCombined[rune(euc_code)] = [2]rune{rune(uni_code1), rune(uni_code2)}
			} else if v, existing := replaceMap[rune(euc_code)]; existing {
				eucJis1st[euc_code-0xA1A1] = v
			} else {
				eucJis1st[euc_code-0xA1A1] = rune(uni_code1)
			}
			break

		case 0x8EA1 <= euc_code && euc_code <= 0x8EDF:
			eucJisHankaku[euc_code-0x8EA1] = rune(uni_code1)
			break

		case 0x8EE0 <= euc_code && euc_code <= 0x8EFE:
			break

		case 0x8FA1A1 <= euc_code && euc_code <= 0x8FFFEF6:
			eucJis2nd[euc_code-0x8FA1A1] = rune(uni_code1)
			break

		default:
			log.Fatal("Invalid code ", s)
		}
	}

	outputFile(eucJis1st[:], eucJisHankaku[:], eucJis2nd[:], eusJis1stCombined, eusJis1stCombinedKeys)
}

func outputFile(eucJis1st []rune, eucJisHankaku []rune, eucJis2nd []rune, eusJis1stCombined map[rune][2]rune, eusJis1stCombinedKeys []rune) {
	fmt.Print("// GENERATED FROM https://x0213.org/codetable/euc-jis-2004-std.txt, DO NOT EDIT\n")
	fmt.Print("package table\n\n")

	// var eucJis1st = [23902]rune{}  //  A1A1(41377) ~ FEFE(65278)
	fmt.Print("var EUC_JIS_1ST_MAP = [...]rune{")
	for idx, v := range eucJis1st {
		if idx%16 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("0x%x, ", v)
	}
	fmt.Print("\n}\n")

	// var eucJisHankaku = [63]rune{} // 8EA1 (36513) ~  8EDF (36575)
	fmt.Print("\nvar EUC_JIS_HANKAKU_MAP = [...]rune{")
	for idx, v := range eucJisHankaku {
		if idx%16 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("0x%x,", v)
	}
	fmt.Print("\n}\n")

	// var eucJis2nd = [23902]rune{}  // 8FA1A1 (9413025) ~ 8FFEFE (9436926)
	fmt.Print("\nvar EUC_JIS_2ND_MAP = [...]rune{")
	for idx, v := range eucJis2nd {
		if idx%16 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("0x%x,", v)
	}
	fmt.Print("\n}\n")

	// var eusJis1stCombined = map[rune][2]rune{}
	sort.Slice(eusJis1stCombinedKeys, func(i, j int) bool {
		return i < j
	})
	fmt.Print("\nvar EUC_JIS_1ST_COMBINED = map[int32][2]rune {\n")
	for _, k := range eusJis1stCombinedKeys {
		v := eusJis1stCombined[k]
		fmt.Printf("0x%x: {0x%x, 0x%x},\n", k, v[0], v[1])
	}
	fmt.Print("}\n")
}

func parseLine(s string, euc_code *int32, uni_code1 *int32, uni_code2 *int32) {
	if _, err := fmt.Sscanf(s, "0x%X	U+%X+%X	", euc_code, uni_code1, uni_code2); err == nil {
		return
	} else if _, err := fmt.Sscanf(s, "0x%X	U+%X	", euc_code, uni_code1); err == nil {
		return
	} else {
		log.Fatal("Invalid line ", s)
	}
}
