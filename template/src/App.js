import Register from "./Register";
import LoginForm from "./LoginForm";
import { useState } from "react";

function App() {
  const [isRegister, setIsRegister] = useState(false);
  return (
    <div className="App">
      <header className="App-header">
        <Register />
        <button onClick={() => setIsRegister(!isRegister)}>Press after registering</button>
        {isRegister && <LoginForm />}
      </header>
    </div>
  );
}

export default App;
