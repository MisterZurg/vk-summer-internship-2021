/*
Ввод-вывод:

• Выражения для вычисления читаются из stdin, по одному выражению на строку.
· Строки разделены \n.

• Ответы пишутся в stdout, по одному ответу на строку.
· Строки разделены \n.
· Ничего, кроме ответа, в stdout не пишется.

• Программа завершает работу после обработки всех данных из stdin.

Формат входных данных:
• Числа, состоящие из римских цифр: I, V, X, L, C, D, M, Z.
· Z обозначает 0.
• Операции: +, -, *, /.
· Унарный минус (-) перед числом обозначает, что оно отрицательное.
• Скобки (, ).
• Между числами, операциями и скобками допустимы пробелы.

Формат ответа:
• В случае успеха — число, состоящее из римских цифр.
• В случае ошибки — строка, начинающаяся с error:.
*/

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type StackRune []rune

func (s *StackRune) Push(v rune) {
	*s = append(*s, v)
}

func (s *StackRune) Pop() rune {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *StackRune) Peek() rune {
	return (*s)[len(*s)-1]
}

type StackInt []int64

func (s *StackInt) Push(v int64) {
	*s = append(*s, v)
}

func (s *StackInt) Pop() int64 {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *StackInt) Peek() int64 {
	return (*s)[len(*s)-1]
}

func main() {
	inputExpression()
}

func inputExpression()  {
	defer func() {
		if r := recover(); r != nil {
			inputExpression()
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println(solutionOfExpression(strings.TrimSpace(s)))
	}
}

//Решение подготовленного выражения в ОПН
func solutionOfExpression(expression string) string {
	expression = strings.ReplaceAll(expression, " ", "")
	// Посмотреть выражение
	// fmt.Println(expression)
	isLegitimate(expression)
	preparedExpression := getAdoptedExpression(expression)
	// fmt.Println("Prep:", preparedExpression)
	reversedPolishNotation := expression2ReversedPolishNotation(preparedExpression)
	arabicNumber := reversedPolishNotation2Answer(reversedPolishNotation)
	// fmt.Println("Arab:",arabicNumber)
	if arabicNumber < 0 {
		cutminus := strconv.FormatInt(arabicNumber, 10)
		ans, _ := strconv.ParseInt(cutminus[1:], 10, 64)
		return "-" + arabicToRoman(ans)
	}
	return arabicToRoman(arabicNumber)
}

func isLegitimate(expression string) {
	numberOfLeftBrackets := strings.Count(expression, "(")
	numberOfRightbrackets := strings.Count(expression, ")")
	if numberOfLeftBrackets != numberOfRightbrackets {
		fmt.Println("error: number of left brackets not equal number of right")
		panic("error: number of left brackets not equal number of right")
	}
	if containsUnallowedSymbol(expression) {
		fmt.Println("error: expressions contains illegal symbol or operand")
		panic("error: expressions contains illegal symbol or operand")
	}
	// Check expressions like (I+II)*
	if 	expression[0] == '*'|| expression[len(expression)-1] == '*' ||
		expression[len(expression)-1] == '-' || // ха-ха  expression[0] == '-'
		expression[0] == '+'|| expression[len(expression)-1] == '+' ||
		expression[0] == '/'|| expression[len(expression)-1] == '/' {
		fmt.Println("error: beginning/ending of expression contains illegal symbol")
		panic("error: beginning/ending of expression contains illegal symbol")
	}
}

func containsUnallowedSymbol(expression string) bool {
	// Внимание: Индусский код в Микрочипе
	for _, symbol := range expression {
		if  symbol != 'I' &&
			symbol != 'V' &&
			symbol != 'X' &&
			symbol != 'L' &&
			symbol != 'C' &&
			symbol != 'D' &&
			symbol != 'M' &&
			symbol != 'Z' &&
			symbol != '+' &&
			symbol != '-' &&
			symbol != '*' &&
			symbol != '/' &&
			symbol != '(' &&
			symbol != ')' {
			return true
		}
	}
	return false
}

func getAdoptedExpression(expression string) string {
	var adoptedExpression string = ""
	for token := 0; token < len(expression); token++ {
		var symbol = rune(expression[token])
		if symbol == '-' {
			if token == 0 { // Чекаем первый символ, явл ли он -
				adoptedExpression += "Z"
			} else if expression[token-1] == '(' {
				adoptedExpression += "Z"
			}
		}
		adoptedExpression += string(symbol)
	}
	return adoptedExpression
}

//Выражение в обратную польскую нотацию
func expression2ReversedPolishNotation(expression string) string {
	var current string
	var stack StackRune

	var currentPriority int
	// Проходимся по выражению посимвольно
	for i := 0; i < len(expression); i++ {
		// Получаем текущий приоритет
		currentPriority = getPriorityOfOperation(rune(expression[i]))
		// Если число
		if currentPriority == 0 {
			current += string(expression[i])
		}
		// Если открывающаяся скоба (
		if currentPriority == 1 {
			stack.Push(rune(expression[i]))
		}
		// Если математический операнд
		if currentPriority > 1 {
			// Разделяем элементы состоящие из более чем 1 цифры
			current += " "
			// Проверим стек на пустоту
			for ; len(stack) > 0; {
				/* Пока он не пустой, записываем в current
				   все элементы, с приоритетом меньше
				   текущего currentPriority
				*/
				if getPriorityOfOperation(rune(stack.Peek())) >= currentPriority {
					current += string(stack.Pop())
				} else {
					break
				}
			}
			stack.Push(rune(expression[i]))
		}
		// Если закрывающаяся скоба )
		if currentPriority == -1 {
			// Разделяем элементы состоящие из более чем 1 цифры
			current += " "
			// Забираем элементы из стека до тех пор,
			// пока не встретим открывающаяся скобку
			for ; getPriorityOfOperation(rune(stack.Peek())) != 1; {
				current += string(stack.Pop())
			}
			stack.Pop()
		}
	}
	for ; len(stack) > 0; {
		current += string(stack.Pop())
	}
	return current
}

//Обратная польская нотация в ответ, как ни странно :/
func reversedPolishNotation2Answer(rpn string) int64 {
	var operand string = ""
	var stack StackInt

	for i := 0; i < len(rpn); i++ {
		// fmt.Println(string(rpn[i]))
		if (rune(rpn[i])) == ' ' {
			continue
		}
		// Если число
		if getPriorityOfOperation(rune(rpn[i])) == 0 {
			// fmt.Println("Число")
			// Собираем все число
			for ; rpn[i] != ' ' && getPriorityOfOperation(rune(rpn[i])) == 0; {
				operand += string(rpn[i])
				i++
				if i == len(rpn) {
					break
				}
			}
			// Для римских чисел
			arabicNumber := romanToArabic(operand)
			stack.Push(arabicNumber)
			// Для арабских чисел in case
			// number, _ := strconv.ParseInt(operand, 10, 64)
			//stack.Push(number)
			operand = ""
		}
		// Если математический операнд
		if getPriorityOfOperation(rune(rpn[i])) > 1 {
			// Забираем из стека два последних числа
			var a int64 = stack.Pop()
			var b int64 = stack.Pop()

			switch rune(rpn[i]) {
			case '+':
				stack.Push(b + a)
			case '-':
				stack.Push(b - a)
			case '*':
				stack.Push(b * a)
			case '/':
				if a == 0 {
					fmt.Println("error: division by zero!")
					panic("error: division by zero!")
				} else {
					ans:= math.Floor(float64(b) /float64(a))
					// fmt.Println(math.Floor(float64(b) /float64(a)), "=" ,ans)
					stack.Push(int64(ans))
				}
			}
		}
	}
	return stack.Pop()
}

func getPriorityOfOperation(token rune) int {
	if token == '*' || token == '/' {
		return 3
	} else if token == '+' || token == '-' {
		return 2
	} else if token == '(' {
		return 1
	} else if token == ')' {
		return -1
	} else {
		return 0 //Приоритет чисел
	}
}

/*
Time Complexity: O(N), where n is the length of the string.
Only one traversal of the string is required.

Space Complexity: O(1).
As no extra space is required.
*/
func romanToArabic(roman string) int64 {
	roman2Arabic := map[string]int64{
		"I": 1,
		"V": 5,
		"X": 10,
		"L": 50,
		"C": 100,
		"D": 500,
		"M": 1000,
		"Z": 0,
	}
	// Initialize previous character and answer
	var p int64 = 0
	var ans int64 = 0

	// Traverse through all characters
	n := utf8.RuneCountInString(roman)
	for i := n - 1; i > -1; i-- {
		// If greater than or equal to previous,
		// add to answer
		if roman2Arabic[string(roman[i])] >= p {
			ans += roman2Arabic[string(roman[i])]
		} else { // If smaller than previous
			ans -= roman2Arabic[string(roman[i])]
		}
		// Update previous
		p = roman2Arabic[string(roman[i])]
	}
	return ans
}
// Time Complexity: O(N);
// Space Complexity: O(1).
func arabicToRoman(arabic int64) string {
	arabicNums := [13]int64{1000,
		900,
		500,
		400,
		100,
		90,
		50,
		40,
		10,
		9,
		5,
		4,
		1,
		//	0,
	}
	romanNums := [13]string{
		"M",
		"CM",
		"D",
		"CD",
		"C",
		"XC",
		"L",
		"XL",
		"X",
		"IX",
		"V",
		"IV",
		"I",
		//	"Z",
	}
	// Для удобства сравнения, не используем тк.
	// ключи не гарантируют нам упорядоченность
	//arabic2Roman := map[int64]string{
	//	1000: "M",
	//	900:  "CM",
	//	500:  "D",
	//	400:  "CD",
	//	100:  "C",
	//	90:   "XC",
	//	50:   "L",
	//	40:   "XL",
	//	10:   "X",
	//	9:    "IX",
	//	5:    "V",
	//	4:    "IV",
	//	1:    "I",
	//	0:    "Z",
	//}

	var ans = ""
	if arabic == 0 {
		ans += "Z"
		return ans
	}
	for ind := range arabicNums {

		var numberOfSymbols = arabic / arabicNums[ind]
		if numberOfSymbols != 0 {
			for i := 0; i < int(numberOfSymbols); i++ {
				ans += romanNums[ind]
			}
		}
		arabic %= arabicNums[ind]
	}
	return ans
}

func testCase() {
	// Тестируем римское -> арабское
	// fmt.Println(romanToArabic("M"))
	// fmt.Println(romanToArabic("XVII")) // to represent 17,
	// fmt.Println(romanToArabic("MCMLIII")) //  for 1953,
	// fmt.Println(romanToArabic("MMMCCCIII")) // for 3303

	//fmt.Println("(XVI + I) + I")
	// Тестируем разделение выражения
	//fmt.Println("((81 * 6) /42+ (3-1))")
	//fmt.Println(expression2ReversedPolishNotation("((81 * 6) /42+ (3-1))"))
	// Тестируем РПН
	// fmt.Println(expression2ReversedPolishNotation("2+2*2")) // 2 2 2*+
	// fmt.Println(expression2ReversedPolishNotation("(2+2)*2")) // 2 2 + 2*
	// fmt.Println(expression2ReversedPolishNotation("222*2")) // 2 2 + 2*

	// Тестируем РПН в ответ
	// fmt.Println(reversedPolishNotation2Answer("2 2 2*+")) // 6

	// Тестируем РПН римскую
	// fmt.Println(expression2ReversedPolishNotation("I+I")) // I I+
	// fmt.Println(expression2ReversedPolishNotation("(XVI+I)+I")) // XVI I + I+

	// Тестируем римскую РПН в ответ
	// fmt.Println(reversedPolishNotation2Answer("I I+"))
	// fmt.Println(reversedPolishNotation2Answer("XVI I + I+"))

	// Тестируем арабское -> римское
	// fmt.Println(arabicToRoman(124))
	// fmt.Println(arabicToRoman(17))// XVII
	// fmt.Println(arabicToRoman(1953))// MCMLIII
	// fmt.Println(arabicToRoman(3303))// MMMCCCIII
}
