package letuchka3;
//вариант 53
import java.util.Scanner;

public class PALParser {
    private static Scanner scanner;
    private static int line = 1, col = 1;

    public static void main(String[] args) {
        scanner = new Scanner(System.in);
        System.out.println("Enter the input string:");
        if (Pal()) {
            System.out.println("Input string is valid.");
        } else {
            System.out.println("Syntax error at (" + line + ", " + col + ")");
        }
    }

    private static boolean Pal() {
        if (scanner.hasNext("IDENT")) {
            System.out.print("IDENT ");
            match("IDENT");
            if (!Pal()) {
                return false;
            }
            if (scanner.hasNext("IDENT")) {
                System.out.print("IDENT");
                match("IDENT");
                return true;
            } else {
                return false;
            }
        } else if (scanner.hasNextInt()) {
            System.out.print("NUMBER ");
            match("NUMBER");
            if (!Pal()) {
                return false;
            }
            if (scanner.hasNextInt()) {
                System.out.print("NUMBER");
                match("NUMBER");
                return true;
            } else {
                return false;
            }
        } else if (scanner.hasNext("STRING")) {
            System.out.print("STRING ");
            match("STRING");
            if (!Pal()) {
                return false;
            }
            if (scanner.hasNext("STRING")) {
                System.out.print("STRING");
                match("STRING");
                return true;
            } else {
                return false;
            }
        } else if (scanner.hasNext("\\(")) {
            System.out.print("( ");
            match("\\(");
            if (!Mid()) {
                return false;
            }
            if (scanner.hasNext("\\)")) {
                System.out.print(") ");
                match("\\)");
                return true;
            } else {
                return false;
            }
        } else {
            return true;
        }
    }

    private static boolean Mid() {
        if (Pal()) {
            System.out.print("<Mid> ::= <Pal> ");
            return Mid();
        } else {
            System.out.print("<Mid> ::= ε");
            return true;
        }
    }

    private static void match(String token) {
        if (scanner.hasNext(token)) {
            String value = scanner.next(token);
            updatePosition(value);
        } else {
            System.out.println("Syntax error at (" + line + ", " + col + ")");
            System.exit(1);
        }
    }

    private static void updatePosition(String value) {
        for (int i = 0; i < value.length(); i++) {
            if (value.charAt(i) == '\n') {
                line++;
                col = 1;
            } else {
                col++;
            }
        }
    }
}