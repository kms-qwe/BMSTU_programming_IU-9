package exam1;

public class Test {
    public static void main(String[] args) {
        Fraction fraction1 = new Fraction(1, 2);
        Fraction fraction2 = new Fraction(3, 4);

        Fraction sum = fraction1.add(fraction2);
        Fraction product = fraction1.multiply(fraction2);

        System.out.println("Sum: " + sum);        
        System.out.println("Product: " + product); 
    }
}
