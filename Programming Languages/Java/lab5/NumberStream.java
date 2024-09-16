package lab5;
import java.util.Arrays;
import java.util.Optional;
import java.util.stream.IntStream;
import java.util.stream.Stream;
public class NumberStream {
    private final int[] numbers;

    public NumberStream(int[] numbers) {
        this.numbers = numbers;
    }

    public Stream<Integer> getDigitStream() {
        return IntStream.of(numbers)
                .mapToObj(number -> String.valueOf(number).chars())
                .flatMapToInt(intStream -> intStream)
                .mapToObj(Character::getNumericValue);
    }
    public Optional<Integer> getMaxDigitSum() {
        return Arrays.stream(numbers)
                .mapToObj(number -> Math.abs(number))
                .map(number -> String.valueOf(number).chars().map(Character::getNumericValue).boxed())
                .map(stream -> stream.reduce(0, Integer::sum)) 
                .max(Integer::compareTo); 
    }

}
