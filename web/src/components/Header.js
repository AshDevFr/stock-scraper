import React from "react";
import styled from 'styled-components';

const HeaderDiv = styled.div`
  background-color: #262626;
  width: 100%;
  margin: 0;
  padding: 10px;
  color: white;

  h2 {
    margin: 0;
    padding: 0;
  }
`

const Header = () => (
  <HeaderDiv>
    <h2>Stock Scraper</h2>
  </HeaderDiv>
);

export default Header;