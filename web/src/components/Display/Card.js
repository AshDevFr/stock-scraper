import React from "react";
import styled from 'styled-components';
import CardTitle from "./CardTitle";

const CardDiv = styled.div`
  border: 1px solid grey;
  border-radius: 5px;
  margin: 5px;
  padding: 5px;
  height: 200px;
`;

const Link = styled.a``;

const Card = ({item}) => (
  <CardDiv>
    <CardTitle item={item} />
    <div>
      Tracked url: <Link href={item.trackedUrl} target="_blank">link</Link>
    </div>
  </CardDiv>
);

export default Card;