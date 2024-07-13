import React, { useState, useEffect } from "react";
import axios from "axios";
import "./App.css";

function App() {
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [expiration, setExpiration] = useState("");
  const [cache, setCache] = useState({});
  const [error, setError] = useState("");

  // Fetch cache data from the server
  const fetchCache = async () => {
    try {
      const response = await axios.get("http://localhost:8080/cache");
      setCache(response.data);
    } catch (err) {
      console.error("Error fetching cache:", err);
    }
  };

  // Handle WebSocket events
  const handleWebSocketEvents = () => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      console.log("WebSocket connection opened");
    };

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      const { event: evt, key, value, expiration } = data;

      if (evt === "set" && value) {
        setCache((prevCache) => ({
          ...prevCache,
          [key]: { value, expiration },
        }));
      } else if (evt === "delete") {
        setCache((prevCache) => {
          const newCache = { ...prevCache };
          delete newCache[key];
          return newCache;
        });
      }
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };

    return socket;
  };

  // Set up WebSocket and fetch initial cache data
  useEffect(() => {
    fetchCache();
    const socket = handleWebSocketEvents();
    return () => socket.close();
  }, []);

  // Handle setting a cache value
  const handleSet = async () => {
    if (!key) {
      setError("Key cannot be empty");
      return;
    }
    if (!value) {
      setError("Value cannot be empty");
      return;
    }

    try {
      await axios.post("http://localhost:8080/cache", {
        key,
        value,
        expiration: parseInt(expiration) || 0,
      });
      setKey("");
      setValue("");
      setExpiration("");
      setError("");
    } catch (err) {
      setError("Error setting cache");
      console.error("Error setting cache:", err);
    }
  };

  // Handle deleting a cache value
  const handleDelete = async (keyToDelete) => {
    setError("");
    try {
      await axios.delete(`http://localhost:8080/cache/${keyToDelete}`);
    } catch (err) {
      setError("Error deleting cache");
      console.error("Error deleting cache:", err);
    }
  };

  return (
    <div className="App">
      <h1>LRU Cache App</h1>
      <div className="form-group">
        <label>Key:</label>
        <input
          type="text"
          value={key}
          onChange={(e) => setKey(e.target.value)}
        />
      </div>
      <div className="form-group">
        <label>Value:</label>
        <input
          type="text"
          value={value}
          onChange={(e) => setValue(e.target.value)}
        />
      </div>
      <div className="form-group">
        <label>Expiration (ms):</label>
        <input
          type="number"
          value={expiration}
          onChange={(e) => setExpiration(e.target.value)}
        />
      </div>
      <button onClick={handleSet}>Set</button>
      {error && <p className="error">{error}</p>}

      {/* Display current cache state */}
      <div className="cache-state">
        <h2>Currently in Cache</h2>
        <ul>
          {Object.keys(cache).length > 0 ? (
            Object.entries(cache).map(([key, { value, expiration }], index) => (
              <li key={index}>
                <strong>{key}:</strong> {value}{" "}
                <em>
                  Expires on: {new Date(expiration / 1e6).toLocaleString()}
                </em>
                <button
                  onClick={() => handleDelete(key)}
                  className="delete-button"
                >
                  Delete
                </button>
              </li>
            ))
          ) : (
            <li>No items in cache</li>
          )}
        </ul>
      </div>
    </div>
  );
}

export default App;
