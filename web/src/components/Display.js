import React from "react";
import {useSelector} from "react-redux";
import styled from 'styled-components';
import Grid from 'styled-components-grid';
import Card from './Display/Card';

const DisplayDiv = styled.div``;

const Display = () => {
  const {config, loaded} = useSelector((state) => state.config);

  if (!loaded) {
    return (<div>Config not loaded</div>);
  }

  return (
    <DisplayDiv>
      <Grid>
        {config.items.map(item => {
          return <Grid.Unit key={item.uuid} size={{desktop: 1 / 5, tablet: 1 / 4}}>
            <Card item={item}/>
          </Grid.Unit>
        })}
      </Grid>
    </DisplayDiv>
  );
}

export default Display;