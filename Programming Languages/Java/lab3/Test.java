package lab3;
import java.util.Arrays;
// var 36 59
public class Test {
    public static void main(String[] args) {
        System.out.println("===============\nVAR36\n");

        MatrixWithRangSort[] matrices = new MatrixWithRangSort[] {
            new MatrixWithRangSort(new int[]{3, 2, 1, 4, 5, 6}, 2, 3),
            new MatrixWithRangSort(new int[]{0, 0, 0, 0, 1, 0, 0, 0, 1}, 3, 3),
        };

        System.out.println("BEFORE SORT");
        for (MatrixWithRangSort matrix : matrices) {
            matrix.printMatrix();
        }

        Arrays.sort(matrices);

        System.out.println("\nAFTER SORT");
        for (MatrixWithRangSort matrix : matrices) {
            matrix.printMatrix();
        }


        System.out.println("===============\nVAR59");

        SimpleFraction[] fractions = new SimpleFraction[10];
        for (int i = 0; i < 10; i++) {
            int numerator = (int)(Math.random() * 10) + 1;
            int denominator = (int)(Math.random() * 10) + 1;
            int sign = (Math.random() > 0.5) ? -1 : 1;
            fractions[i] = new SimpleFraction(sign*numerator, denominator);
        }
        
        System.out.println("BEFORE SORT");
        for (SimpleFraction fraction : fractions) {
            System.out.print(fraction + " ");
        }
        System.err.println();

        Arrays.sort(fractions);

        System.out.println("AFTER SORT");
        for (SimpleFraction fraction : fractions) {
            System.out.print(fraction + " ");
        }
        System.err.println();
    }
}
