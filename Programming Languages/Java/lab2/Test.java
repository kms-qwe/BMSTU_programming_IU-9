package lab2;

public class Test {
    public static void main(String[] args){
        String[] names = {"Alice", "Bob", "Charlie", "Karl"};
        int[] rus = {10000, 30, 100, 20};
        int[] math = {20, 40, 40, 10};
        int[] inf = {1, 0, 1, 0};
        boolean[] at = {true, false, true, true};
        Students students = new Students(names, rus, math, inf, at);
        System.out.println(students.toSring());
        Students best = students.Best(2);
        System.out.println(best.toSring());
    }

}
