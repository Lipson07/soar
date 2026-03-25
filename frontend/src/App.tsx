import { Route, Routes } from "react-router-dom";
import "./App.scss";
import MeetPage from "./Components/MeetPage/MeetPage";
import Login from "./Components/Form/Login/Login";

function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<MeetPage />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </>
  );
}

export default App;
