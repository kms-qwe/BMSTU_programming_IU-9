package lab5;

import java.util.Comparator;
import java.util.Optional;
import java.util.stream.IntStream;
import java.util.stream.Stream;

public class PalindromSubstringStream {
    private String s;
    private int k;

    public PalindromSubstringStream(String s, int k){
        this.s = s;
        this.k = k;
    }

    public Stream<String> getSubstringPalindromStream() {
        return IntStream.rangeClosed(0, s.length())
                .boxed()
                .flatMap(start -> IntStream.rangeClosed(start + 1, s.length())
                        .mapToObj(end -> s.substring(start, end)))
                .filter(this::isPalindrome);
    }

    private  boolean isPalindrome(String s) {
        return s.length() == k && s.equals(new StringBuilder(s).reverse().toString());
    }

    public Optional<Integer> getFirstIndex() {
        Optional<Integer> ind = Optional.empty();
        Optional<String> firstSubPol = getSubstringPalindromStream()
                .sorted(Comparator.naturalOrder())
                .findFirst(); 
        if (firstSubPol.isPresent()) {
            String pal = firstSubPol.get();
            ind = Optional.of(s.indexOf(pal));
        }
        return ind;
    }

}
