package lab4;
// var 34
// Строка, составленная из маленьких латинских
// букв, с итератором по всем подстрокам, в которых
// количество гласных букв совпадает с количеством
// согласных букв.

// var 56
// Строка, составленная из маленьких латинских
// букв, с итератором по максимальным
// «правильным» подстрокам. В «правильной»
// подстроке гласные и согласные буквы
// чередуются. 
public class Test {
     public static void main(String[] args) {
        String str = "fffaccbe";
        EqSubstring sub = new EqSubstring(str);
        for (String s : sub) {
            System.out.println(s);
        }
        AlterSubstring sub2 = new AlterSubstring(str);
        for (String s : sub2) {
            System.out.println(s);
        }

    }
}