import React, { createContext, useContext, useState } from 'react';
import {
  // fetchGroups,
  addGroup,
  editGroup,
  removeGroup
} from '../services/groupsService';

export const GroupsContext = createContext();

export const GroupsProvider = ({ children }) => {
  const [groups, setGroups] = useState([]);
  // const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // const getGroups = async () => {
  //   try {
  //     setLoading(true);
  //     const data = await fetchGroups();
  //     setGroups(data.groups || []);
  //     return { success: true };
  //   } catch (err) {
  //     setError(err);
  //     return { success: false, error: err };
  //   } finally {
  //     setLoading(false);
  //   }
  // };

  // useEffect(() => {
  //   getGroups();
  // }, []);

  const createGroup = async groupData => {
    try {
      const newGroup = await addGroup(groupData);
      setGroups(prevGropus => [...prevGropus, newGroup]);
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const updateGroup = async (id, groupData) => {
    try {
      await editGroup(id, groupData);

      // Update the local state with the latest group data
      setGroups(prevGropus =>
        prevGropus.map(group =>
          group.id === id ? { ...group, ...groupData } : group
        )
      );
      return { success: true };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  const deleteGroup = async id => {
    try {
      await removeGroup(id);
      setGroups(prevGropus => prevGropus.filter(group => group.id !== id));
      return { success: true, id };
    } catch (err) {
      setError(err);
      return { success: false, error: err };
    }
  };

  return (
    <GroupsContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        groups,
        // getGroups,
        // loading,
        error,
        createGroup,
        updateGroup,
        deleteGroup
      }}
    >
      {children}
    </GroupsContext.Provider>
  );
};

export const useGroups = () => {
  return useContext(GroupsContext);
};
