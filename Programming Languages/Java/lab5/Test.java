package lab5;


import java.util.Comparator;
import java.util.Map;
import java.util.stream.Collectors;

//27 Последовательность целых чисел с операциями:
// 1. порождение потока цифр десятичного
// представления чисел (например, для
// последовательности 10, 0, 123, 5 поток должен
// содержать цифры 1, 0, 0, 1, 2, 3, 5);
// 2. поиск максимальной суммы цифр числа
// последовательности.
// Проверить работу первой операции нужно путём
// подсчёта количеств каждой из десятичных цифр в
// потоке.

// 46 Последовательность символов Unicode с
// операциями:
// 1. порождение потока строк, представляющих все
// подпоследовательности длины k, являющиеся
// палиндромами;
// 2. поиск индекса первой буквы
// лексикографически наименьшего палиндрома
// длины k.
// Проверить работу первой операции нужно путём
// ранжирования палиндромов их потока по
// количеству содержащихся в них различных букв. 
public class Test {
    public static void main(String[] args) {
        int[] numbers = {10, 0, 123, 5};
        NumberStream numberStream = new NumberStream(numbers);

        System.out.println("Digits Stream:");
        numberStream.getDigitStream().forEach(System.out::println);

        System.out.println("Max Digit Sum: " + numberStream.getMaxDigitSum().orElse(0));

        Map<Integer, Long> digitFrequency = numberStream.getDigitStream()
                .collect(Collectors.groupingBy(
                        digit -> digit, Collectors.counting()));

        System.out.println("Digit Frequency:");
        digitFrequency.forEach((digit, frequency) ->
                System.out.println(digit + ": " + frequency));


        String s = "CABABABBB";
        int k = 3;
        PalindromSubstringStream PSStream = new PalindromSubstringStream(s, k);

        System.out.println("\nSubstrings Palindrom len " + k + " sorted by number of unique characters:");

        PSStream.getSubstringPalindromStream()
                .sorted(Comparator.comparingInt(Test::countUniqueCharacters))
                .forEach(System.out::println);

        System.out.println("Index Of First Palindrom lexicographically:");
        System.out.println(PSStream.getFirstIndex().orElse(-1));

      
    }
    public static int countUniqueCharacters(String s) {
        return (int) s.chars().distinct().count();
    }
}
