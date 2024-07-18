import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [nombre, setNombre] = useState('');
  const [cantidad, setCantidad] = useState('');
  const [unidad, setUnidad] = useState('');
  const [insumos, setInsumos] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchInsumos = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/stock');
        setInsumos(response.data);
      } catch (error) {
        console.error('Error al obtener insumos:', error);
        console.log('Error details:', error.response);
      }
    };
    fetchInsumos();
  }, []);

  const handleEliminar = async (id) => {
    try {
      await axios.delete(`http://localhost:8080/api/stock/${id}`);
      setInsumos(insumos.filter((insumo) => insumo.id !== id));
    } catch (error) {
      setError('Error al eliminar el insumo: ' + error.message);
    }
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setLoading(true);

    try {
      const existingInsumo = insumos.find(i => i.nombre === nombre);

      if (existingInsumo) {
        await axios.put(`http://localhost:8080/api/stock/${existingInsumo.id}`, {
          cantidad: existingInsumo.cantidad + parseInt(cantidad)
        });
      } else {
        await axios.post('http://localhost:8080/api/stock', {
          nombre: nombre,
          cantidad: parseInt(cantidad),
          unidad: unidad
        });
      }

      const response = await axios.get('http://localhost:8080/api/stock');
      setInsumos(response.data);

      setNombre('');
      setCantidad('');
      setUnidad('');
    } catch (error) {
      setError('Error al agregar o actualizar el insumo: ' + error.message);
    } finally {
      setLoading(false);
    }
  };
  const handleDescSubmit = async (event) => {
    event.preventDefault();
    setLoading(true);

    try {
      const existingInsumo = insumos.find(i => i.nombre === nombre);

      if (existingInsumo) {
        await axios.put(`http://localhost:8080/api/stock/${existingInsumo.id}`, {
          cantidad: existingInsumo.cantidad - parseInt(cantidad)
        });
      } else {
        await axios.post('http://localhost:8080/api/stock', {
          nombre: nombre,
          cantidad: parseInt(cantidad),
          unidad: unidad
        });
      }

      const response = await axios.get('http://localhost:8080/api/stock');
      setInsumos(response.data);

      setNombre('');
      setCantidad('');
      setUnidad('');
    } catch (error) {
      setError('Error al descontar el insumo: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Panadería del Valle</h1>
      </header>
      <h2>Insumos:</h2>
      <table>
        <thead>
          <tr>
            <th>Insumos</th>
            <th>Cantidad</th>
            <th>Unidad</th>
            <th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {insumos.map((insumo) => (
            <tr key={insumo.id}>
              <td>{insumo.nombre}</td>
              <td>{insumo.cantidad}</td>
              <td>{insumo.unidad}</td>
              <td>
                <button onClick={() => handleEliminar(insumo.id)}>Eliminar</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <h2>Actualizar Insumo:</h2>
      <h3>Acordate de poner el nombre igual a como se encuentra en la tabla.</h3>
      <form onSubmit={handleSubmit}>
        <label>
          Nombre:
          <input
            type="text"
            value={nombre}
            onChange={(e) => setNombre(e.target.value)}
            required
          />
        </label>
        <br />
        <label>
          Cantidad:
          <input
            type="number"
            value={cantidad}
            onChange={(e) => setCantidad(e.target.value)}
            required
          />
        </label>
        <br />
        <br />
        <label>
          Unidad:
          <input
            type="text"
            value={unidad}
            onChange={(e) => setUnidad(e.target.value)}
            required
          />
        </label>
        <br />
        <button type="submit" disabled={loading}>{loading ? 'Enviando...' : 'Agregar'}</button>
      </form>
      <form onSubmit={handleDescSubmit}>
        <button type="desc" disabled={loading}>{loading ? 'Enviando...' : 'Descontar'}</button>
      </form>
      {error && <p className="error">{error}</p>}
    </div>
  );
}

export default App;
