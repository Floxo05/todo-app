import AuthHelper from "../auth/Auth";
import {Todo} from "../../components/Todo/TodoList";

class TodoService {
    static async getTodos() {
        const response = await fetch(process.env.REACT_APP_API + '/auth/todos', {
            headers: {
                'Authorization': 'Bearer ' + AuthHelper.getToken() || '',
            }
        });
        const data = await response.json();
        return this.returnData(data);
    }

    static async addTodo(todo: Todo) {
        const response = await fetch( process.env.REACT_APP_API +'/auth/todo/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + AuthHelper.getToken() || '',
            },
            body: JSON.stringify({title: todo.title}),
        });
        const data = await response.json();
        return this.returnData(data);
    }

    static async updateTodoStatus(todo: Todo) {
        const response = await fetch(process.env.REACT_APP_API + `/auth/todo/${todo.id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + AuthHelper.getToken() || '',
            },
            body: JSON.stringify(todo),
        });
        const data = await response.json();
        return this.returnData(data);
    }

    static async deleteTodoById(id: number) {
        const response = await fetch(process.env.REACT_APP_API + `/auth/todo/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + AuthHelper.getToken() || '',
            },
        });

        const data = await response.json();

        return this.returnData(data);
    }

    private static returnData(data: any) {
        // if data has attribute error, throw error with data.error
        if (data.error) {
            throw new Error(data.error);
        }
        return data;
    }
}

export default TodoService;