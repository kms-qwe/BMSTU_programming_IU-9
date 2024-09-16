package lab2;

public class Students {

    private Student[] students;

    public Students(String[] names, int[] rus, int[] math, int[] inf, boolean[] attestat) {
        if (names.length == rus.length && rus.length == math.length && math.length == inf.length && inf.length == attestat.length) {
            students = new Student[names.length];
            for (int i = 0; i < names.length; i++) {
                students[i] = new Student(names[i], rus[i], math[i], inf[i], attestat[i]);
            }
        }
    }
    public Students Best(int endIndex){
        int len = 0;
        for (int i = 0; i < this.students.length; i++) {
            if (this.students[i].getAttestat()) {
                len += 1;
            }
        }
        Students res = new Students(new String[len], new int[len], new int[len], new int[len], new boolean[len]);
        res.students = new Student[len];
        int top = 0;
        for (int i = 0; i < this.students.length; i++) {
            if (this.students[i].getAttestat()) {
                res.students[top] = this.students[i];
                top += 1;
            }
        }
        for (int i = 0; i < res.students.length;i ++) {
            int n = res.students.length;
            for (int j = 0; j < n - i -1; j++) {
                if (res.students[j].compareTo(res.students[j + 1]) < 0) {
                    Student temp = res.students[j];
                    res.students[j] = res.students[j + 1];
                    res.students[j + 1] = temp;
                }
            }
        }
        if (endIndex > len) {
            endIndex = len;
        }
        Students res2 = new Students(new String[endIndex], new int[endIndex], new int[endIndex], new int[endIndex], new boolean[endIndex]);
        res2.students = new Student[endIndex];
        for (int i = 0; i < endIndex; i++) {
            res2.students[i] = res.students[i];
        }
        return res2;
    }
    public String toSring(){
        String res = "";
        for (int i = 0; i < this.students.length; i++) {
            res += this.students[i].toString();
            res += "\n";
        }
        return res;
    }
}