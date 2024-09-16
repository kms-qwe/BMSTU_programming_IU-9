package lab4;

import java.util.Iterator;
public class EqSubstring implements Iterable<String> {
    private String str;

    public EqSubstring (String str) {
        this.str = str.toLowerCase();
    }

    public Iterator<String> iterator() {
        return new VowelConsonantIterator();
    }
    private class VowelConsonantIterator implements Iterator<String> {
        private int start = 0;
        private int end = 0;

        public boolean hasNext() {
            return start < str.length();
        }


        public String next() {
            if (end  >= str.length()) {
                start += 1;
                end = start;
            }
            if (!hasNext()) {
                return "END ITERATOR";
            }
            String substring = str.substring(start, end + 1);
            int vowels = 0;
            int consonants = 0;

            end += 1;
            for (int i = 0; i < substring.length(); i++) {
                char ch = substring.charAt(i);
                if (isVowel(ch)) {
                    vowels++;
                } else {
                    consonants++;
                }
            }
          
            if (vowels == consonants) {
                return substring.toString();
            }
            return next();
        }

        private boolean isVowel(char ch) {
            return "aeiou".indexOf(ch) != -1;
        }
    }
}
