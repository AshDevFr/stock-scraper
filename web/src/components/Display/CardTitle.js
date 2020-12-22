import React from "react";
import styled from 'styled-components';

const CardTitleDiv = styled.div`
  font-size: large;
  font-weight: bold;
  width: 100%;
  padding-bottom: 10px;
`

const CardTitle = ({item}) => (
  <CardTitleDiv>
    {item.type.toUpperCase()} {item.id}
  </CardTitleDiv>
);

export default CardTitle;