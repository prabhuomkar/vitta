import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Checkbox,
  Card,
  CardBody
} from '@chakra-ui/react';
import { CheckCircleIcon } from '@chakra-ui/icons';
import { useState } from 'react';

function ConversionTable() {
  const [selectedAll, setSelectedAll] = useState(false);
  const [selectedRows, setSelectedRows] = useState({});
  const [clearedStatus, setClearedStatus] = useState(false);

  const data = [
    {
      id: 1,
      name: 'XYZ Bank',
      payee: 'Zomato',
      category: 'Food',
      credit: 0,
      debit: 100,
      cleared: true
    },
    {
      id: 2,
      name: 'XYZ Bank',
      payee: 'Amazon',
      category: 'Fashion',
      credit: 0,
      debit: 200,
      cleared: true
    },
    {
      id: 3,
      name: 'XYZ Bank',
      payee: 'Swiggy',
      category: 'Food',
      credit: 0,
      debit: 300,
      cleared: true
    }
  ];

  const handleSelectAll = () => {
    const newValue = !selectedAll;
    setSelectedAll(newValue);
    setSelectedRows(
      newValue
        ? data.reduce((acc, item) => ({ ...acc, [item.id]: true }), {})
        : {}
    );
  };

  const handleRowSelect = id => {
    setSelectedRows(prev => {
      const newSelection = { ...prev, [id]: !prev[id] };
      if (!newSelection[id]) delete newSelection[id];
      setSelectedAll(Object.keys(newSelection).length === data.length);
      return newSelection;
    });
  };

  const toggleCleared = id => {
    setClearedStatus(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  return (
    <Card boxShadow="none" width="auto">
      <CardBody>
        <TableContainer>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>
                  <Checkbox
                    isChecked={selectedAll}
                    onChange={handleSelectAll}
                  />
                </Th>
                <Th>#</Th>
                <Th>Name</Th>
                <Th>Payee</Th>
                <Th>Category</Th>
                <Th isNumeric>Credit</Th>
                <Th isNumeric>Debit</Th>
                <Th>Cleared</Th>
              </Tr>
            </Thead>
            <Tbody>
              {data.map((item, index) => (
                <Tr key={item.id}>
                  <Td>
                    <Checkbox
                      isChecked={!!selectedRows[item.id]}
                      onChange={() => handleRowSelect(item.id)}
                    />
                  </Td>
                  <Td>{index + 1}</Td>
                  <Td>{item.name}</Td>
                  <Td>{item.payee}</Td>
                  <Td>{item.category}</Td>
                  <Td isNumeric>{item.credit}</Td>
                  <Td isNumeric>{item.debit}</Td>
                  <Td>
                    <CheckCircleIcon
                      color={clearedStatus[item.id] ? 'green.500' : 'gray.500'}
                      cursor="pointer"
                      onClick={() => toggleCleared(item.id)}
                    />
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </TableContainer>
      </CardBody>
    </Card>
  );
}

export default ConversionTable;
