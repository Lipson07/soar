import { Route, Routes } from "react-router-dom";
import "./App.scss";
import MeetPage from "./Components/Pages/MeetPage/MeetPage";
import Login from "./Components/Form/Login/Login";
import Register from "./Components/Form/Register/Register";
import MainPage from "./Components/Pages/MainPage/MainPage";

function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<MeetPage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/main" element={<MainPage />} />
      </Routes>
    </>
  );
}

export default App;
