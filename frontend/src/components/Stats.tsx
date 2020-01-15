import React, { useEffect, useState } from "react";
import {
    Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
    Paper,
} from "@material-ui/core";

const Stats: React.FC = () => {
    const [counts, setCounts] = useState([]);
    useEffect(() => {
        fetch("/api/stats").then((res: Response) => {
            if (!res.ok) {
                throw new Error(res.statusText);
            }
            return res.json();
        }).then((results) => {
            setCounts(results);
        }).catch((err: Error) => {
            window.console.error(err.message);
        });
    }, []);
    const rows = counts.map((v: any, index: number) => {
        return (
          <TableRow key={index}>
            <TableCell component="th" scope="row">Size {v["size"]}</TableCell>
            <TableCell>{v["status_ready"]}</TableCell>
            <TableCell>{v["status_ng"]}</TableCell>
            <TableCell>{v["status_pending"]}</TableCell>
            <TableCell>{v["status_ok"]}</TableCell>
            <TableCell>{v["status_predicted"]}</TableCell>
          </TableRow>
        );
    });
    return (
      <React.Fragment>
        <h2>Stats</h2>
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell></TableCell>
                <TableCell>Ready</TableCell>
                <TableCell>NG</TableCell>
                <TableCell>Pending</TableCell>
                <TableCell>OK</TableCell>
                <TableCell>Predicted</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{rows}</TableBody>
          </Table>
        </TableContainer>
      </React.Fragment>
    );
};

export default Stats;
