import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

// Конфигурация API
const API_URL = 'https://beagle-mighty-terribly.ngrok-free.app/api/v1';

const api = axios.create({
  baseURL: API_URL,
});

// Интерцептор для добавления JWT токена в заголовки запросов
api.interceptors.request.use(
  async (config) => {
    const token = await AsyncStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Авторизация пользователя
export const login = async (email, password) => {
  const response = await api.post('/login', { email, password });
  if (response.data.token) {
    await AsyncStorage.setItem('token', response.data.token);
  }
  return response.data;
};

// Регистрация пользователя
export const register = async (email, password) => {
  const response = await api.post('/regin', { email, password });
  console.log('Register response:', response.data);
  return response.data;
};

// Выход из системы
export const logout = async () => {
  await AsyncStorage.removeItem('token');
};

// Получение списка чатов пользователя
export const getChats = async (userId) => {
  try {
    const response = await api.get(`/chats/getchats/${userId}`);
    return response.data;
  } catch (error) {
    console.error('Error fetching chats:', error);
    throw error;
  }
};

// Получение сообщений по чату
export const getMessages = async (chatId) => {
  try {
    const response = await api.get(`/chats/getmessages/${chatId}`);
    return response.data;
  } catch (error) {
    console.error('Ошибка получения сообщений:', error);
    throw error;
  }
};
