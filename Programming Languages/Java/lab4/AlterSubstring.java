package lab4;

import java.util.Iterator;
public class AlterSubstring implements Iterable<String> {
    private String str;

    public AlterSubstring (String str) {
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
            end += 1;
            for (int i = 0; i < substring.length() - 1; i++) {
                char chI = substring.charAt(i);     
                char ch2I = substring.charAt(i+1);
                if (isVowel(chI) && !isVowel(ch2I) || !isVowel(chI) && isVowel(ch2I)) {
                    continue;
                } 
                return next();

            }
            return substring;
        }

        private boolean isVowel(char ch) {
            return "aeiou".indexOf(ch) != -1;
        }
    }
}