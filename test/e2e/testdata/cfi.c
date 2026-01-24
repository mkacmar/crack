typedef int (*func_ptr)(int);

int add_one(int x) {
    return x + 1;
}

int multiply_two(int x) {
    return x * 2;
}

int call_indirect(func_ptr fn, int x) {
    return fn(x);
}

int main(void) {
    func_ptr fn = add_one;
    int result = call_indirect(fn, 5);
    fn = multiply_two;
    result += call_indirect(fn, 3);
    return result > 0 ? 0 : 1;
}
