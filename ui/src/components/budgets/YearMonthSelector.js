import React from 'react';
import {
  Box,
  Flex,
  Heading,
  Button,
  Tab,
  Tabs,
  TabList
} from '@chakra-ui/react';

const YearMonthSelector = ({
  selectedYear,
  selectedMonth,
  months,
  currentYear,
  primaryColor,
  setSelectedYear,
  setSelectedMonth,
  setIsAddGroupOpen
}) => {
  const getSurroundingMonths = () => {
    const surroundingMonths = [];
    const currentDate = new Date();
    const currentMonth = currentDate.getMonth();

    // Loop from -2 to 2 to get 2 previous, current, and 2 next months
    // eslint-disable-next-line no-plusplus
    for (let i = -2; i <= 2; i++) {
      const date = new Date(currentYear, currentMonth + i);
      surroundingMonths.push({
        name: months[date.getMonth()],
        year: date.getFullYear(),
        value: date.getMonth() + 1
      });
    }

    return surroundingMonths;
  };

  const surroundingMonths = getSurroundingMonths();

  // Find the index of the selected month and year in surroundingMonths array
  const selectedIndex = surroundingMonths.findIndex(
    month => month.value === selectedMonth && month.year === selectedYear
  );

  return (
    <Box mb={4}>
      <Flex
        justifyContent="space-between"
        alignItems="center"
        flexDirection={{ base: 'column', md: 'row' }}
      >
        <Heading mb={{ base: 4, md: 0 }} mr={2}>
          {selectedYear}
        </Heading>
        <Flex
          alignItems="center"
          flexDirection={{ base: 'column', md: 'row' }}
          width="100%"
          overflow="scroll"
        >
          <Tabs
            variant="soft-rounded"
            size="sm"
            index={selectedIndex}
            onChange={index => {
              const { year, value } = surroundingMonths[index];
              setSelectedYear(year);
              setSelectedMonth(value);
            }}
          >
            <TabList>
              {surroundingMonths.map(({ name }, index) => (
                <Tab
                  // eslint-disable-next-line react/no-array-index-key
                  key={index}
                  color="black"
                  _selected={{ color: 'white', bg: `${primaryColor}` }}
                  p="0.4rem 0.8rem"
                  ml={1}
                >
                  {name}
                </Tab>
              ))}
            </TabList>
          </Tabs>
        </Flex>
        <Button
          colorScheme={primaryColor}
          bg={primaryColor}
          padding="1rem 1.4rem"
          onClick={() => setIsAddGroupOpen(true)}
          width={{ base: '100%', md: 'auto' }}
          ml={{ base: 0, md: 2 }}
          mt={{ base: 4, md: 0 }}
        >
          Add Group
        </Button>
      </Flex>
    </Box>
  );
};

export default YearMonthSelector;
