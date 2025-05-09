import http

import streamlit as st
import requests

# Базовый URL для Go-сервера
API_BASE_URL = "http://localhost:8080"

st.set_page_config(
    page_title="The Sweets Of Life",
    page_icon="🍰",
)

# Инициализация состояния
if "jwt_token" not in st.session_state:
    st.session_state["jwt_token"] = None

# Функция для выполнения запросов с JWT токеном
def api_request(method, endpoint, **kwargs):
    headers = kwargs.pop("headers", {})
    if st.session_state["jwt_token"]:
        headers["Authorization"] = f"Bearer {st.session_state['jwt_token']}"
    response = requests.request(method, f"{API_BASE_URL}{endpoint}", headers=headers, **kwargs)
    return response

# Функции для API-запросов
def sign_up(fullname, username, email, password, phone_number):
    response = requests.post(
        f"{API_BASE_URL}/auth/sign-up",
        json={
            "fullname": fullname,
            "username": username,
            "email": email,
            "password_hash": password,
            "phone_number": phone_number,
        },
    )
    return response

def sign_in(username, password):
    response = requests.post(f"{API_BASE_URL}/auth/sign-in", json={"username": username, "password": password})
    return response

def make_order(user_id, delivery, cakes, payment_method):
    data = {
        "user_id": user_id,
        "delivery": delivery,
        "cakes": cakes,
        "payment_method": payment_method
    }
    response = api_request("POST", "/api/make-order", json=data)
    return response

def view_orders():
    response = api_request("GET", "/api/view-orders")
    return response

def cancel_order(order_id):
    data = {"order_id": order_id}
    response = api_request("POST", "/api/delete-order", json=data)
    return response

def get_cakes():
    response = api_request("GET", "/api/cakes")
    return response

def get_delivery_points():
    response = api_request("GET", "/api/delivery-points")
    return response

def admin_sign_in(username, password):
    response = requests.post(f"{API_BASE_URL}/adm/sign-in", json={"username": username, "password": password})
    return response

def admin_add(username, password):
    response = requests.post(f"{API_BASE_URL}/adm/add-admin", json={"username": username, "password_hash": password})
    return response

def authorization_page():
    st.title("Авторизация")
    auth_action = st.radio("Выберите действие", ["Войти", "Войти как администратор", "Зарегистрироваться"])

    username = st.text_input("Никнейм")
    password = st.text_input("Пароль", type="password")

    if auth_action == "Войти":
        if st.button("Войти"):
            response = sign_in(username, password)
            if response.status_code == 200:
                result = response.json()
                st.success("Вы успешно вошли в систему!")
                st.session_state["jwt_token"] = result.get("token")
                st.session_state["role"] = "user"
                st.rerun()
            else:
                st.error("Ошибка авторизации")
    elif auth_action == "Войти как администратор":
        if st.button("Войти как администратор"):
            response = admin_sign_in(username, password)
            if response.status_code == 200:
                result = response.json()
                st.success("Вы вошли как администратор!")
                st.session_state["jwt_token"] = result.get("token")
                st.session_state["role"] = "admin"
                st.rerun()
            else:
                st.error("Ошибка авторизации администратора")
    elif auth_action == "Зарегистрироваться":
        fullname = st.text_input("ФИО")
        email = st.text_input("Электронная почта")
        phone_number = st.text_input("Номер телефона")
        if st.button("Зарегистрироваться"):
            response = sign_up(fullname, username, email, password, phone_number)
            if response.status_code == 200:
                st.success("Регистрация успешна!")
            else:
                st.error("Ошибка регистрации")

def manage_users_page():
    st.title("Управление пользователями")

    # Получение списка пользователей
    response = api_request("GET", "/adm/manage-users/users")
    if response.status_code == 200:
        users = response.json()
        for user in users:
            st.subheader(f"Пользователь {user['fullname']} ({user['email']})")
            st.text(f"Никнейм: {user['username']}")
            st.text(f"Телефон: {user['phone_number']}")

    else:
        st.error("Ошибка загрузки списка пользователей")


def manage_cakes_page():
    st.title("Управление тортами")

    # Получение списка тортов
    response = api_request("GET", "/adm/manage-cakes/cakes")
    if response.status_code == 200:
        cakes = response.json()
        for cake in cakes:
            st.subheader(f"{cake['description']} - {cake['price']} $")
            st.text(f"ID: {cake['id']}, Вес: {cake['weight']} г")

            # Удалить торт
            if st.button(f"Удалить {cake['description']}", key=cake['id']): 
                delete_response = api_request("POST", f"/adm/manage-cakes/remove-cake", json={"id": cake["id"]})
                if delete_response.status_code == 200:
                    st.success(f"{cake['description']} удален!")
                else:
                    st.error("Ошибка удаления торта")

        # Добавление нового торта
        st.subheader("Добавить новый торт")
        new_description = st.text_input("Название")
        new_price = st.number_input("Цена", min_value=0.0, step=0.5)
        new_weight = st.number_input("Вес (г)", min_value=0, step=50)
        new_full_description = st.text_input("Описание")
        if st.button("Добавить торт"):
            add_response = api_request("POST", "/adm/manage-cakes/add-cake", json={
                "description": new_description,
                "price": int(new_price),
                "weight": int(new_weight),
                "full_description": new_full_description
            })
            if add_response.status_code == 200:
                st.success("Торт успешно добавлен!")
            else:
                st.error("Ошибка добавления торта")


# def database_management_page():
#     st.title("Управление базой данных")
#
#     # Бэкап базы данных
#     if st.button("Создать бэкап базы данных"):
#         backup_response = api_request("POST", "/adm/database/backup")
#         if backup_response.status_code == 200:
#             st.success("Бэкап успешно создан!")
#         else:
#             st.error("Ошибка создания бэкапа")
#
#     # Восстановление базы данных
#     if st.button("Восстановить базу данных"):
#         recovery_response = api_request("POST", "/adm/database/recovery")
#         if recovery_response.status_code == 200:
#             st.success("База данных успешно восстановлена!")
#         else:
#             st.error("Ошибка восстановления базы данных")

def add_admin_page():
    st.title("Добавление нового администратора ")

    username = st.text_input("Никнейм")
    password = st.text_input("Пароль", type="password")

    if st.button("Зарегистрировать админа"):
        response = admin_add(username, password)
        if response.status_code == http.HTTPStatus.OK:
            st.success("Новый администратор успешно добавлен!")
        else:
            st.error("Ошибка при добавлении нового администратора")


def main():
    st.sidebar.title("Меню")

    # Проверка наличия JWT токена и отображение соответствующего интерфейса
    if st.session_state["jwt_token"]:
        if st.session_state.get("role") == "admin":
            menu = st.sidebar.radio("Навигация", ["Управление пользователями", "Управление тортами", "Добавить администратора", "Выйти"])
            if menu == "Управление пользователями":
                manage_users_page()
            elif menu == "Управление тортами":
                manage_cakes_page()
            # elif menu == "Управление базой данных":
            #     database_management_page()
            elif menu == "Добавить администратора":
                add_admin_page()
            elif menu == "Выйти":
                st.session_state["jwt_token"] = None
                st.session_state["role"] = None
                st.success("Вы вышли из системы!")
                st.rerun()
        else:
            # Маршрутизация на основе состояния
            if st.session_state.get("current_page") == "cake_detail":
                cake_id = st.session_state.get("current_cake_id")
                if cake_id:
                    cake_detail_page(cake_id)
            else:
                menu = st.sidebar.radio("Навигация", ["Каталог", "МАИ заказы", "Сделать заказ", "Выйти"])
                if menu == "Каталог":
                    catalog_page()
                elif menu == "МАИ заказы":
                    orders_page()
                elif menu == "Сделать заказ":
                    create_order_page()
                elif menu == "Выйти":
                    st.session_state["jwt_token"] = None
                    st.success("Вы вышли из системы!")
                    st.rerun()  # Перезагрузка страницы после выхода
    else:
        menu = st.sidebar.radio("Навигация", ["Авторизация"])
        if menu == "Авторизация":
            authorization_page()

def update_order(order_id, payment_method):
    data = {"order_id": order_id, "payment_method": payment_method}
    response = api_request("POST", "/api/change-order", json=data)
    return response

# Страница заказов
def orders_page():
    st.title("Ваши заказы")
    response = view_orders()

    if response.status_code == 200:
        try:
            orders_data = response.json()["Orders"]
            order_avg_cost = response.json()["AvgCost"]
            if not orders_data:
                st.text("Список заказов пока пуст! Купите уже что-нибудь :]")
                return
            for order_data in orders_data:
                order = order_data["Ord"]
                cakes = order_data["Cakes"]

                # Выводим информацию о заказе
                st.subheader(f"Заказ #{str(order['id'])}")
                st.text(f"Дата заказа: {order['time']}")
                st.text(f"Статус: {order['order_status']}")
                st.text(f"Способ оплаты: {order['payment_method']}")
                st.text(f"Стоимость: {order['cost']} $.")

                for cake in cakes:
                    st.text(f"Торт: {cake['description']} ({cake['price']} $, {cake['weight']} г)")

                # Форма для изменения способа оплаты
                with st.expander(f"Изменить способ оплаты для заказа #{order['id']}"):
                    new_payment_method = st.radio(
                        "Выберите новый способ оплаты",
                        ["Card", "Cash", "Online Payment"],
                        index=["Card", "Cash", "Online Payment"].index(order["payment_method"]),
                        key=f"payment_method_{order['id']}"  # Уникальный ключ с использованием ID заказа
                    )
                    if st.button(f"Изменить способ оплаты для заказа #{order['id']}"):
                        update_response = update_order(order['id'], new_payment_method)
                        if update_response.status_code == 200:
                            st.success(f"Способ оплаты для заказа #{order['id']} успешно обновлен!")
                        else:
                            st.error("Ошибка при обновлении способа оплаты")

                # Кнопка для отмены заказа
                if st.button(f"Отменить заказ #{order['id']}"):
                    cancel_response = cancel_order(order['id'])
                    if cancel_response.status_code == 200:
                        st.warning(f"Заказ #{order['id']} отменен!")
                    else:
                        st.error("Ошибка отмены заказа")
            st.text(f"Средняя стоимость заказов: {order_avg_cost} $")
        except KeyError:
            st.error("Ответ сервера не содержит ключа 'Orders' или неправильный формат данных.")
    else:
        st.warning("Ошибка загрузки заказов")


def create_order_page():
    st.title("Сделать заказ")

    # Получаем список всех тортов
    cakes_response = get_cakes()
    if cakes_response.status_code == 200:
        cakes = cakes_response.json()
        selected_cakes = []

        st.subheader("Выберите торты и количество")
        for cake in cakes:
            # Отображение информации о торте
            st.markdown(f"**{cake['description']}**")
            st.text(f"Цена: {cake['price']} $ | Вес: {cake['weight']} г")

            # Поле для выбора количества
            quantity = st.number_input(
                f"Количество {cake['description']}",
                min_value=0,
                max_value=100,
                step=1,
                key=f"quantity_{cake['id']}"
            )
            if quantity > 0:
                for _ in range(quantity):
                    selected_cakes.append({"cake": cake, "quantity": 1})

        if not selected_cakes:
            st.warning("Пожалуйста, выберите хотя бы один торт для заказа.")
            return
    else:
        st.error("Ошибка загрузки списка тортов.")
        return

    # Получаем список пунктов доставки
    delivery_response = get_delivery_points()
    if delivery_response.status_code == 200:
        delivery_points = delivery_response.json()
        delivery_point = st.selectbox("Выберите пункт доставки", [point['address'] for point in delivery_points])
    else:
        st.error("Ошибка загрузки пунктов доставки.")
        return

    # Способ оплаты
    payment_method = st.radio("Выберите способ оплаты", ["Card", "Cash", "Online Payment"])

    if st.button("Оформить заказ"):
        # Считаем общую стоимость и вес
        total_cost = sum(cake['cake']['price'] * cake['quantity'] for cake in selected_cakes)
        total_weight = sum(cake['cake']['weight'] * cake['quantity'] for cake in selected_cakes)

        # Формируем данные для заказа
        selected_cake_items = [
            {"id": cake["cake"]["id"], "description": cake["cake"]["description"],
             "price": cake["cake"]["price"], "weight": cake["cake"]["weight"],
             "quantity": cake["quantity"]}
            for cake in selected_cakes
        ]
        delivery_point_id = next(point['id'] for point in delivery_points if point['address'] == delivery_point)

        order_data = {
            "user_id": 3,  # Нужно заменить на динамический ID пользователя
            "delivery": {
                "point_id": delivery_point_id,
                "cost": total_cost,
                "status": "pending",
                "weight": total_weight,
            },
            "cakes": selected_cake_items,
            "payment_method": payment_method
        }

        order_response = make_order(order_data["user_id"], order_data["delivery"], order_data["cakes"], order_data["payment_method"])

        if order_response.status_code == 200:
            st.success("Заказ успешно оформлен!")
        else:
            st.error("Ошибка оформления заказа.")

def catalog_page():
    st.title("Каталог тортов")
    response = get_cakes()
    if response.status_code == 200:
        cakes = response.json()

        st.markdown(
            """
            <style>
                .cake-card {
                    display: flex;
                    flex-direction: row;
                    align-items: center;
                    justify-content: space-between;
                    border: 1px solid #ddd;
                    border-radius: 10px;
                    padding: 15px;
                    margin-bottom: 20px;
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
                }
                .cake-image {
                    width: 120px;
                    height: 120px;
                    object-fit: cover;
                    border-radius: 10px;
                    margin-right: 15px;
                }
                .cake-info {
                    flex: 1;
                    display: flex;
                    flex-direction: column;
                    justify-content: center;
                }
                .cake-price {
                    font-weight: bold;
                    margin-top: 10px;
                }
                .cake-actions {
                    display: flex;
                    align-items: center;
                    justify-content: center;
                }
            </style>
            """,
            unsafe_allow_html=True
        )

        for cake in cakes:
            image_url = cake.get('image_url')
            if not image_url:
                print("image_url is incorrect:", image_url)
                image_url = "https://img.freepik.com/free-photo/chocolate-cake-with-blueberry-cream_140725-10903.jpg"

            st.markdown(
                f"""
                <div class="cake-card">
                    <img src="{image_url}" class="cake-image" />
                    <div class="cake-info">
                        <h4>{cake["description"]}</h4>
                        <p class="cake-price">Цена: {cake['price']} $</p>
                    </div>
                </div>
                """,
                unsafe_allow_html=True
            )

            # Добавляем кнопку выбора с обработкой состояния
            if st.button(f"Выбрать {cake['description']}", key=cake["id"]):
                st.session_state["current_cake_id"] = cake["id"]
                st.session_state["current_page"] = "cake_detail"
                st.rerun()
    else:
        st.warning("Ошибка загрузки каталога")


def cake_detail_page(cake_id):
    st.title("Детальная информация о торте")
    response = api_request("GET", f"/api/cakes/{cake_id}")  # Запрос данных торта по ID
    if response.status_code == 200:
        cake = response.json()

        image_url = cake.get('image_url')
        if not image_url:
            print("image_url is incorrect:", image_url)
            image_url = "https://img.freepik.com/free-photo/chocolate-cake-with-blueberry-cream_140725-10903.jpg"

        st.image(image_url, use_container_width=True)
        st.subheader(cake["description"])
        st.text(f"Цена: {cake['price']} $")
        st.text(f"Вес: {cake['weight']} г")
        st.text("Описание:")
        st.write(cake["full_description"])

        # Кнопка "Назад"
        if st.button("Назад в каталог"):
            st.session_state["current_page"] = "catalog"  # Перенаправляем пользователя на каталог
            st.rerun()  # Перезагружаем страницу, чтобы отобразить каталог
    else:
        st.error("Ошибка загрузки информации о торте")

if __name__ == "__main__":
    main()
