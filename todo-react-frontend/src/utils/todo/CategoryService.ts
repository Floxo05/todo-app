import AuthHelper from "../auth/Auth";


class CategoryService {
    static async getCategories() {
        const response = await fetch(process.env.REACT_APP_API + '/auth/categories', {
            headers: {
                'Authorization': 'Bearer ' + AuthHelper.getToken() || ''
            }
        });
        return await response.json();
    }
}

export default CategoryService;