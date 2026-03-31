const storage = {
  getItem: (key: string) => {
    try {
      const value = localStorage.getItem(key);
      return Promise.resolve(value ? JSON.parse(value) : null);
    } catch (error) {
      console.warn("Error getting item from localStorage:", error);
      return Promise.resolve(null);
    }
  },
  setItem: (key: string, value: any) => {
    try {
      localStorage.setItem(key, JSON.stringify(value));
      return Promise.resolve();
    } catch (error) {
      console.warn("Error setting item to localStorage:", error);
      return Promise.resolve();
    }
  },
  removeItem: (key: string) => {
    try {
      localStorage.removeItem(key);
      return Promise.resolve();
    } catch (error) {
      console.warn("Error removing item from localStorage:", error);
      return Promise.resolve();
    }
  },
};

export default storage;
