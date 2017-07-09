package arqsoft.authentication.model;

import javax.persistence.*;

/**
 * Based on javergarav on 15/02/2017.
 * Created by heartsTeam
 */

@Entity
@Table(name = "users")
@NamedQueries({@NamedQuery(name = User.FIND_ALL, query = "SELECT u FROM User u")})
public class User {

    public static final String FIND_ALL = "User.findAll";

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private long id;

    private String name;
    private String email;
	private enum roll {OWNER,ADMINSYSTEM,SELLER,CUSTOMER,ADMIN};

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

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public enum getRoll() {
        return roll;
    }

    public void setRoll(enum roll) {
        this.roll = roll;
    }
}