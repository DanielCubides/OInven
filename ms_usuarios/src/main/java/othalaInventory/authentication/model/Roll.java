package othalaInventory.authentication.model;

import javax.persistence.*;
import java.util.*;

/**
 * Based on javergarav 
 * Created by heartsTeam 8/7/2017   
 */

@Entity
@Table(name = "rolls")
@NamedQueries({@NamedQuery(name = Roll.FIND_ALL, query = "SELECT r FROM Roll r")})
public class Roll {

    public static final String FIND_ALL = "Roll.findAll";

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private long id;

    private String name;


    private ArrayList<User> users;

    public long getId() {
        return id;
    }

    public void setId(long id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @ManyToMany(mappedBy = "user")
    public ArrayList<User> getUsers() {
        return users;
    }

    public void setUsers(ArrayList<User> users) {
        this.users = users;
    }
}