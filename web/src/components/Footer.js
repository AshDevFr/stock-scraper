import React from "react";
import styled from 'styled-components';

const FooterContainer = styled.div`
  margin-top: 1rem;
  padding: 1rem;
  position: fixed;
  bottom: 0;
  left: 0;
  
  width: 100%;
  max-height: 200px;
  overflow-x: hidden;
  overflow-y: auto;
  
  color: white;
  background-color: black;
`;

const Footer = ({children}) => {
  return (
    <FooterContainer>
      {children}
    </FooterContainer>
  );
}

export default Footer;