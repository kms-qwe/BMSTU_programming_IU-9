// package lab4;

// import java.util.Iterator;

// public class SuffixList implements Iterable<String> {
//     private StringBuilder s;

//     public SuffixList(StringBuilder s) {
//         this.s = s;
//     }

//     public Iterator<String> iterator() {
//         return new SuffixIterator();
//     }

//     private class SuffixIterator implements Iterator<String> {
//         private int pos;

//         public SuffixIterator() {
//             pos = 0;
//         }

//         public boolean hasNext() {
//             return pos < s.length();
//         }

//         public String next() {
//             return s.substring(pos++, s.length());
//         }
//     }
// }
