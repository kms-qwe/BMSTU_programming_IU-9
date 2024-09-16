package lab2;

public class Student{
    private String name;
    private int rus, math, inf;
    private boolean attestat;

    public Student(String name, int rus, int math, int inf, boolean attestat) {
        this.name = name;
        this.rus = rus;
        this.math = math;
        this.inf = inf;
        this.attestat = attestat;
    }
    public boolean getAttestat() {
        return this.attestat;
    }
    public int sum() {
        return this.rus + this.math + this.inf;
    }
    public int compareTo(Student other) {
        return this.sum() - other.sum();
    }
    public String toString() {
        return name + " " + rus + " " + math + " " + inf + " " + attestat;
    }
}