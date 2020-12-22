import React, {useEffect} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import './App.css';
import Header from "./components/Header.js";

import {fetchConfigAction} from "./actions/configSlice";
import WebSocketProvider from "./components/WebSocketProvider";
import Display from "./components/Display";
import Footer from "./components/Footer";
import Terminal from "./components/Terminal";

const App = () => {
  const dispatch = useDispatch();
  const {loaded} = useSelector((state) => state.config);

  useEffect(() => {
    if (!loaded) {
      dispatch(fetchConfigAction());
    }
  }, [loaded, dispatch]);

  return (
    <WebSocketProvider>
      <Header/>
      <Display/>
      <Footer>
        <Terminal/>
      </Footer>
    </WebSocketProvider>
  );
}

export default App;
