import { Route, Routes } from "react-router-dom";
import "./App.scss";
import MeetPage from "./Components/MeetPage/MeetPage";
import Login from "./Components/Form/Login/Login";
import Register from "./Components/Form/Register/Register";

function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<MeetPage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Routes>
    </>
  );
}

export default App;
